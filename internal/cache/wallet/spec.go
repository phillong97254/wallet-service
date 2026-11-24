package wallet

import (
	"context"

	"github.com/phillong97254/wallet-service/internal/domain"
)

// Cache giữ bản sao số dư ví để giảm tải DB. Chỉ dùng cho đường đọc; ghi số dư
// luôn đọc từ DB (FOR UPDATE), không bao giờ tin cache.
type Cache interface {
	// Get trả ErrMiss nếu ví chưa có trong cache.
	Get(ctx context.Context, adminID string) (domain.Wallet, error)
	Set(ctx context.Context, w domain.Wallet) error
	Delete(ctx context.Context, adminID string) error
}
