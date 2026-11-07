package wallet

import (
	"context"

	walletSvc "github.com/phil/wallet-service/internal/service/wallet"
	walletv1 "github.com/phil/wallet-service/pkg/gen/wallet/v1"
)

func (i *impl) AdjustBalance(ctx context.Context, req *walletv1.AdjustBalanceRequest) (*walletv1.Wallet, error) {
	out, err := i.walletSvc.AdjustBalance(ctx, walletSvc.AdjustBalanceInput{
		AdminID: req.GetAdminId(),
		Amount:  req.GetAmount(),
		Note:    req.GetNote(),
	})
	if err != nil {
		return nil, i.toStatus(err)
	}
	return toProtoWallet(out), nil
}
