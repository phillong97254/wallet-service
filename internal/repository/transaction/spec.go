package transaction

import (
	"context"

	"github.com/phillong97254/wallet-service/internal/domain"
)

type Repository interface {
	// ListByAdminID trả về limit giao dịch mới nhất (limit <= 0: tất cả), thứ tự cũ -> mới.
	ListByAdminID(ctx context.Context, adminID string, limit int) ([]domain.Transaction, error)
	// Create gán ID cho giao dịch rồi lưu.
	Create(ctx context.Context, tx domain.Transaction) (domain.Transaction, error)
}
