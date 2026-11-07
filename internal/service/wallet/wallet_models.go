package wallet

import (
	"time"

	"github.com/phil/wallet-service/internal/domain"
)

type WalletOutput struct {
	AdminID   string
	Currency  string
	Balance   int64
	UpdatedAt time.Time
}

func toWalletOutput(w domain.Wallet) WalletOutput {
	return WalletOutput{AdminID: w.AdminID, Currency: w.Currency, Balance: w.Balance, UpdatedAt: w.UpdatedAt}
}
