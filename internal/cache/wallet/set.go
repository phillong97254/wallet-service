package wallet

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/phillong97254/wallet-service/internal/domain"
)

func (i *impl) Set(ctx context.Context, w domain.Wallet) error {
	raw, err := json.Marshal(w)
	if err != nil {
		return fmt.Errorf("encode wallet: %w", err)
	}
	if err := i.client.Set(ctx, key(w.AdminID), raw, i.ttl).Err(); err != nil {
		return fmt.Errorf("redis set: %w", err)
	}
	return nil
}
