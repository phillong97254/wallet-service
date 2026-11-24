package wallet

import (
	"context"
	"fmt"
)

func (i *impl) Delete(ctx context.Context, adminID string) error {
	if err := i.client.Del(ctx, key(adminID)).Err(); err != nil {
		return fmt.Errorf("redis del: %w", err)
	}
	return nil
}
