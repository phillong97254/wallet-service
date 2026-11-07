package transaction

import (
	"context"
	"slices"

	"github.com/phillong97254/wallet-service/internal/domain"
)

func (i *impl) ListByAdminID(_ context.Context, adminID string, limit int) ([]domain.Transaction, error) {
	i.store.Mu.RLock()
	defer i.store.Mu.RUnlock()
	items := i.store.Transactions[adminID]
	if limit > 0 && limit < len(items) {
		items = items[len(items)-limit:]
	}
	return slices.Clone(items), nil
}
