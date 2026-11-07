package wallet

import (
	"context"
	"fmt"
)

func (i *impl) ListTransactions(ctx context.Context, input ListTransactionsInput) (ListTransactionsOutput, error) {
	if err := canAccess(ctx, input.AdminID); err != nil {
		return ListTransactionsOutput{}, err
	}
	txs, err := i.repo.Transaction().ListByAdminID(ctx, input.AdminID, input.Limit)
	if err != nil {
		return ListTransactionsOutput{}, fmt.Errorf("list transactions: %w", err)
	}
	out := ListTransactionsOutput{Items: make([]TransactionOutput, 0, len(txs))}
	for _, t := range txs {
		out.Items = append(out.Items, TransactionOutput{
			ID: t.ID, Amount: t.Amount, Note: t.Note, CreatedBy: t.CreatedBy, CreatedAt: t.CreatedAt,
		})
	}
	return out, nil
}
