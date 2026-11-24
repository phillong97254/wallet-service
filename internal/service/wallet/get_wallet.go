package wallet

import (
	"context"
	"errors"
	"fmt"

	walletCache "github.com/phillong97254/wallet-service/internal/cache/wallet"
	"github.com/phillong97254/wallet-service/internal/domain"
	walletRepo "github.com/phillong97254/wallet-service/internal/repository/wallet"
)

// GetWallet đọc cache trước, miss thì đọc DB rồi nạp lại cache. Quyền kiểm tra trước khi chạm cache.
// Lỗi cache không làm hỏng request: chỉ log rồi đọc DB.
func (i *impl) GetWallet(ctx context.Context, input GetWalletInput) (WalletOutput, error) {
	if err := canAccess(ctx, input.AdminID); err != nil {
		return WalletOutput{}, err
	}
	if w, ok := i.getCached(ctx, input.AdminID); ok {
		return toWalletOutput(w), nil
	}

	w, err := i.repo.Wallet().GetByAdminID(ctx, input.AdminID)
	if errors.Is(err, walletRepo.ErrNotFound) {
		return WalletOutput{}, fmt.Errorf("wallet of %q not found: %w", input.AdminID, ErrNotFound)
	}
	if err != nil {
		return WalletOutput{}, fmt.Errorf("get wallet: %w", err)
	}
	if err := i.cache.Set(ctx, w); err != nil {
		i.logger.Warn("cache wallet failed", "admin_id", input.AdminID, "err", err)
	}
	return toWalletOutput(w), nil
}

func (i *impl) getCached(ctx context.Context, adminID string) (domain.Wallet, bool) {
	w, err := i.cache.Get(ctx, adminID)
	if err == nil {
		return w, true
	}
	if !errors.Is(err, walletCache.ErrMiss) {
		i.logger.Warn("read wallet cache failed", "admin_id", adminID, "err", err)
	}
	return domain.Wallet{}, false
}
