package routes

import (
	"google.golang.org/grpc"

	v1 "github.com/phillong97254/wallet-service/internal/handler/grpc/v1"
	walletv1 "github.com/phillong97254/wallet-service/pkg/gen/wallet/v1"
)

func V1Routes(srv *grpc.Server, h v1.Handler) {
	walletv1.RegisterWalletServiceServer(srv, h.GetWalletHandler())
}
