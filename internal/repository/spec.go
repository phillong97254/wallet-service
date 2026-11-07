package repository

import (
	"context"

	"github.com/phillong97254/wallet-service/internal/repository/transaction"
	"github.com/phillong97254/wallet-service/internal/repository/wallet"
)

type Registry interface {
	Wallet() wallet.Repository
	Transaction() transaction.Repository
	// Tx chạy fn trong 1 transaction. Bản in-memory: tuần tự hoá các thao tác ghi
	// để đọc-kiểm tra-ghi số dư không bị race (giống SELECT ... FOR UPDATE).
	Tx(ctx context.Context, fn func(tx Registry) error) error
}
