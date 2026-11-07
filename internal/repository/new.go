package repository

import (
	"context"
	"sync"

	"gorm.io/gorm"

	"github.com/phillong97254/wallet-service/internal/repository/transaction"
	"github.com/phillong97254/wallet-service/internal/repository/wallet"
)

type impl struct {
	db *gorm.DB

	wallet     wallet.Repository
	walletOnce sync.Once

	transaction     transaction.Repository
	transactionOnce sync.Once
}

func New(db *gorm.DB) Registry {
	return &impl{db: db}
}

func (i *impl) Wallet() wallet.Repository {
	i.walletOnce.Do(func() { i.wallet = wallet.New(i.db) })
	return i.wallet
}

func (i *impl) Transaction() transaction.Repository {
	i.transactionOnce.Do(func() { i.transaction = transaction.New(i.db) })
	return i.transaction
}

// Tx chạy fn trong 1 DB transaction: registry truyền vào fn dùng *gorm.DB của transaction,
// nên mọi repository (lấy qua getter) đều chạy trong transaction đó.
func (i *impl) Tx(ctx context.Context, fn func(tx Registry) error) error {
	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(New(tx))
	})
}
