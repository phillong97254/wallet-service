package wallet

import "time"

type ListTransactionsInput struct {
	AdminID string
	Limit   int
}

type ListTransactionsOutput struct {
	Items []TransactionOutput
}

type TransactionOutput struct {
	ID        string
	Amount    int64
	Note      string
	CreatedBy string
	CreatedAt time.Time
}
