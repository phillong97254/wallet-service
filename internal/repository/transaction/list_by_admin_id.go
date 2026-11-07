package transaction

import (
	"context"
	"slices"

	"github.com/phillong97254/wallet-service/internal/domain"
)

func (i *impl) ListByAdminID(ctx context.Context, adminID string, limit int) ([]domain.Transaction, error) {
	var items []domain.Transaction
	q := i.db.WithContext(ctx).Where("admin_id = ?", adminID).Order("seq DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&items).Error; err != nil {
		return nil, err
	}
	slices.Reverse(items) // cũ -> mới
	return items, nil
}
