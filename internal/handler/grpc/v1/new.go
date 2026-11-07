package v1

import (
	"log/slog"

	walletHandler "github.com/phillong97254/wallet-service/internal/handler/grpc/v1/wallet"
	walletSvc "github.com/phillong97254/wallet-service/internal/service/wallet"
	walletv1 "github.com/phillong97254/wallet-service/pkg/gen/wallet/v1"
)

type impl struct {
	wallet walletv1.WalletServiceServer
}

func New(logger *slog.Logger, walletService walletSvc.Service) Handler {
	return &impl{wallet: walletHandler.New(logger, walletService)}
}

func (i *impl) GetWalletHandler() walletv1.WalletServiceServer { return i.wallet }
