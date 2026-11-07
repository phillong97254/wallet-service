package repository

import (
	"context"
	"sync"

	"github.com/phillong97254/wallet-service/internal/repository/memstore"
	"github.com/phillong97254/wallet-service/internal/repository/transaction"
	"github.com/phillong97254/wallet-service/internal/repository/wallet"
)

type impl struct {
	store *memstore.Store
	txMu  *sync.Mutex

	wallet     wallet.Repository
	walletOnce sync.Once

	transaction     transaction.Repository
	transactionOnce sync.Once
}

func New(store *memstore.Store) Registry {
	return &impl{store: store, txMu: &sync.Mutex{}}
}

func (i *impl) Wallet() wallet.Repository {
	i.walletOnce.Do(func() { i.wallet = wallet.New(i.store) })
	return i.wallet
}

func (i *impl) Transaction() transaction.Repository {
	i.transactionOnce.Do(func() { i.transaction = transaction.New(i.store) })
	return i.transaction
}

func (i *impl) Tx(_ context.Context, fn func(tx Registry) error) error {
	i.txMu.Lock()
	defer i.txMu.Unlock()
	// Clone registry cho transaction. Thêm repository mới thì nhớ thêm vào đây
	// (với DB thật: truyền *gorm.DB của transaction vào từng repository).
	return fn(&impl{store: i.store, txMu: i.txMu})
}
