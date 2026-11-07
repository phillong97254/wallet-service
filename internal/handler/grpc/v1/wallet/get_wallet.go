package wallet

import (
	"context"

	walletSvc "github.com/phil/wallet-service/internal/service/wallet"
	walletv1 "github.com/phil/wallet-service/pkg/gen/wallet/v1"
)

func (i *impl) GetWallet(ctx context.Context, req *walletv1.GetWalletRequest) (*walletv1.Wallet, error) {
	out, err := i.walletSvc.GetWallet(ctx, walletSvc.GetWalletInput{AdminID: req.GetAdminId()})
	if err != nil {
		return nil, i.toStatus(err)
	}
	return toProtoWallet(out), nil
}
