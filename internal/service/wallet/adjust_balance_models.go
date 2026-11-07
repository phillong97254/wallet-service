package wallet

type AdjustBalanceInput struct {
	AdminID string
	Amount  int64 // âm = trừ tiền
	Note    string
}
