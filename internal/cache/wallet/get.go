package wallet

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/phillong97254/wallet-service/internal/domain"
)

func (i *impl) Get(ctx context.Context, adminID string) (domain.Wallet, error) {
	raw, err := i.client.Get(ctx, key(adminID)).Bytes()
	if errors.Is(err, redis.Nil) {
		return domain.Wallet{}, ErrMiss
	}
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("redis get: %w", err)
	}
	var w domain.Wallet
	if err := json.Unmarshal(raw, &w); err != nil {
		return domain.Wallet{}, fmt.Errorf("decode cached wallet: %w", err)
	}
	return w, nil
}
