package wallet

import (
	"context"

	"github.com/phillong97254/wallet-service/internal/domain"
)

func (i *impl) Save(_ context.Context, w domain.Wallet) error {
	i.store.Mu.Lock()
	defer i.store.Mu.Unlock()
	i.store.Wallets[w.AdminID] = w
	return nil
}
