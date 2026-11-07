package v1

import walletv1 "github.com/phil/wallet-service/pkg/gen/wallet/v1"

type Handler interface {
	GetWalletHandler() walletv1.WalletServiceServer
}
