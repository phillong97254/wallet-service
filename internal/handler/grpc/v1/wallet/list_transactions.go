package wallet

import (
	"context"

	walletSvc "github.com/phil/wallet-service/internal/service/wallet"
	walletv1 "github.com/phil/wallet-service/pkg/gen/wallet/v1"
)

func (i *impl) ListTransactions(ctx context.Context, req *walletv1.ListTransactionsRequest) (*walletv1.ListTransactionsResponse, error) {
	out, err := i.walletSvc.ListTransactions(ctx, walletSvc.ListTransactionsInput{
		AdminID: req.GetAdminId(),
		Limit:   int(req.GetLimit()),
	})
	if err != nil {
		return nil, i.toStatus(err)
	}
	resp := &walletv1.ListTransactionsResponse{Items: make([]*walletv1.Transaction, 0, len(out.Items))}
	for _, t := range out.Items {
		resp.Items = append(resp.Items, &walletv1.Transaction{
			Id: t.ID, Amount: t.Amount, Note: t.Note, CreatedBy: t.CreatedBy, CreatedAtUnix: t.CreatedAt.Unix(),
		})
	}
	return resp, nil
}
