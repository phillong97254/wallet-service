package transaction

import (
	"context"
	"fmt"

	"github.com/phillong97254/wallet-service/internal/domain"
)

// Create gán ID tx-<n> theo từng ví rồi lưu. Gọi trong Tx sau khi đã khoá dòng ví
// (GetByAdminIDForUpdate) để 2 giao dịch đồng thời không lấy trùng số.
func (i *impl) Create(ctx context.Context, tx domain.Transaction) (domain.Transaction, error) {
	var n int64
	if err := i.db.WithContext(ctx).Model(&domain.Transaction{}).Where("admin_id = ?", tx.AdminID).Count(&n).Error; err != nil {
		return domain.Transaction{}, err
	}
	tx.ID = fmt.Sprintf("tx-%d", n+1)
	if err := i.db.WithContext(ctx).Create(&tx).Error; err != nil {
		return domain.Transaction{}, err
	}
	return tx, nil
}
