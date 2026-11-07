package wallet

import (
	"context"
	"errors"
	"fmt"
	"time"

	auth "github.com/phillong97254/go-lib/authkit"
	"github.com/phillong97254/wallet-service/internal/domain"
	"github.com/phillong97254/wallet-service/internal/repository"
	walletRepo "github.com/phillong97254/wallet-service/internal/repository/wallet"
)

// AdjustBalance cộng/trừ số dư và ghi giao dịch (audit created_by) trong cùng 1 transaction.
// Quyền (wallet:write + superadmin) đã được interceptor kiểm tra.
func (i *impl) AdjustBalance(ctx context.Context, input AdjustBalanceInput) (WalletOutput, error) {
	if input.Amount == 0 {
		return WalletOutput{}, ErrInvalidAmount
	}
	c, ok := auth.ClaimsFrom(ctx)
	if !ok {
		return WalletOutput{}, ErrUnauthenticated
	}

	var updated domain.Wallet
	err := i.repo.Tx(ctx, func(tx repository.Registry) error {
		w, err := tx.Wallet().GetByAdminID(ctx, input.AdminID)
		if errors.Is(err, walletRepo.ErrNotFound) {
			return fmt.Errorf("wallet of %q not found: %w", input.AdminID, ErrNotFound)
		}
		if err != nil {
			return fmt.Errorf("get wallet: %w", err)
		}
		if w.Balance+input.Amount < 0 {
			return ErrInsufficientBalance
		}

		now := time.Now().UTC()
		w.Balance += input.Amount
		w.UpdatedAt = now
		if err := tx.Wallet().Save(ctx, w); err != nil {
			return fmt.Errorf("save wallet: %w", err)
		}
		if _, err := tx.Transaction().Create(ctx, domain.Transaction{
			AdminID: input.AdminID, Amount: input.Amount, Note: input.Note,
			CreatedBy: c.Subject, CreatedAt: now,
		}); err != nil {
			return fmt.Errorf("create transaction: %w", err)
		}
		updated = w
		return nil
	})
	if err != nil {
		return WalletOutput{}, err
	}

	i.logger.Info("balance adjusted", "admin_id", input.AdminID, "amount", input.Amount, "by", c.Subject)
	return toWalletOutput(updated), nil
}
