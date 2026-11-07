package repository

import (
	"context"

	"github.com/phillong97254/wallet-service/internal/repository/transaction"
	"github.com/phillong97254/wallet-service/internal/repository/wallet"
)

type Registry interface {
	Wallet() wallet.Repository
	Transaction() transaction.Repository
	// Tx chạy fn trong 1 DB transaction (commit nếu fn trả nil, rollback nếu lỗi).
	Tx(ctx context.Context, fn func(tx Registry) error) error
}
