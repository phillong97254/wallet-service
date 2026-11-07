package wallet

import "context"

type Service interface {
	GetWallet(ctx context.Context, input GetWalletInput) (WalletOutput, error)
	ListTransactions(ctx context.Context, input ListTransactionsInput) (ListTransactionsOutput, error)
	AdjustBalance(ctx context.Context, input AdjustBalanceInput) (WalletOutput, error)
}
