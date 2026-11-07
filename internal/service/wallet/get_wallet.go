package wallet

import (
	"context"
	"errors"
	"fmt"

	walletRepo "github.com/phil/wallet-service/internal/repository/wallet"
)

func (i *impl) GetWallet(ctx context.Context, input GetWalletInput) (WalletOutput, error) {
	if err := canAccess(ctx, input.AdminID); err != nil {
		return WalletOutput{}, err
	}
	w, err := i.repo.Wallet().GetByAdminID(ctx, input.AdminID)
	if errors.Is(err, walletRepo.ErrNotFound) {
		return WalletOutput{}, fmt.Errorf("wallet of %q not found: %w", input.AdminID, ErrNotFound)
	}
	if err != nil {
		return WalletOutput{}, fmt.Errorf("get wallet: %w", err)
	}
	return toWalletOutput(w), nil
}
