package wallet

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/phillong97254/wallet-service/internal/domain"
)

func (i *impl) GetByAdminID(ctx context.Context, adminID string) (domain.Wallet, error) {
	var w domain.Wallet
	err := i.db.WithContext(ctx).Where("admin_id = ?", adminID).Take(&w).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.Wallet{}, ErrNotFound
	}
	return w, err
}
