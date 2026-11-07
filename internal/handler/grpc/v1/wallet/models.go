package wallet

import (
	walletSvc "github.com/phillong97254/wallet-service/internal/service/wallet"
	walletv1 "github.com/phillong97254/wallet-service/pkg/gen/wallet/v1"
)

func toProtoWallet(out walletSvc.WalletOutput) *walletv1.Wallet {
	return &walletv1.Wallet{
		AdminId:       out.AdminID,
		Currency:      out.Currency,
		Balance:       out.Balance,
		UpdatedAtUnix: out.UpdatedAt.Unix(),
	}
}
