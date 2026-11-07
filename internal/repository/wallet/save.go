package wallet

import (
	"context"

	"github.com/phillong97254/wallet-service/internal/domain"
)

func (i *impl) Save(ctx context.Context, w domain.Wallet) error {
	return i.db.WithContext(ctx).Model(&domain.Wallet{}).Where("admin_id = ?", w.AdminID).
		Updates(map[string]any{"balance": w.Balance, "updated_at": w.UpdatedAt}).Error
}
