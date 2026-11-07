package wallet

import (
	"context"

	"github.com/phillong97254/wallet-service/internal/domain"
)

type Repository interface {
	GetByAdminID(ctx context.Context, adminID string) (domain.Wallet, error)
	Save(ctx context.Context, w domain.Wallet) error
}
