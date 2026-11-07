package wallet

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	walletSvc "github.com/phillong97254/wallet-service/internal/service/wallet"
)

// toStatus đổi lỗi nghiệp vụ của service sang gRPC status.
func (i *impl) toStatus(err error) error {
	switch {
	case errors.Is(err, walletSvc.ErrUnauthenticated):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, walletSvc.ErrForbidden):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, walletSvc.ErrNotFound):
		// message dạng `wallet of "x" not found`
		msg := err.Error()
		if u := errors.Unwrap(err); u != nil {
			msg = msg[:len(msg)-len(": "+u.Error())]
		}
		return status.Error(codes.NotFound, msg)
	case errors.Is(err, walletSvc.ErrInvalidAmount):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, walletSvc.ErrInsufficientBalance):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		i.logger.Error("internal error", "error", err)
		return status.Error(codes.Internal, "internal error")
	}
}
