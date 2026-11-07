package transaction

import (
	"context"
	"fmt"

	"github.com/phil/wallet-service/internal/domain"
)

func (i *impl) Create(_ context.Context, tx domain.Transaction) (domain.Transaction, error) {
	i.store.Mu.Lock()
	defer i.store.Mu.Unlock()
	tx.ID = fmt.Sprintf("tx-%d", len(i.store.Transactions[tx.AdminID])+1)
	i.store.Transactions[tx.AdminID] = append(i.store.Transactions[tx.AdminID], tx)
	return tx, nil
}
