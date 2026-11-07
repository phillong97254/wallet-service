package wallet

import (
	"context"

	"github.com/phil/wallet-service/internal/domain"
)

func (i *impl) GetByAdminID(_ context.Context, adminID string) (domain.Wallet, error) {
	i.store.Mu.RLock()
	defer i.store.Mu.RUnlock()
	w, ok := i.store.Wallets[adminID]
	if !ok {
		return domain.Wallet{}, ErrNotFound
	}
	return w, nil
}
