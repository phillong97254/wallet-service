package wallet

import (
	"log/slog"

	walletSvc "github.com/phillong97254/wallet-service/internal/service/wallet"
	walletv1 "github.com/phillong97254/wallet-service/pkg/gen/wallet/v1"
)

// impl hiện thực walletv1.WalletServiceServer (interface sinh từ proto đóng vai trò spec.go).
type impl struct {
	walletv1.UnimplementedWalletServiceServer
	logger    *slog.Logger
	walletSvc walletSvc.Service
}

func New(logger *slog.Logger, walletService walletSvc.Service) walletv1.WalletServiceServer {
	return &impl{logger: logger.With("handler", "wallet"), walletSvc: walletService}
}
