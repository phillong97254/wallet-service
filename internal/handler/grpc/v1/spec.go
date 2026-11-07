package v1

import walletv1 "github.com/phillong97254/wallet-service/pkg/gen/wallet/v1"

type Handler interface {
	GetWalletHandler() walletv1.WalletServiceServer
}
