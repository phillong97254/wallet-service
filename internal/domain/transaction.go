package domain

import "time"

// Transaction là một lần thay đổi số dư. CreatedBy là audit: ai đã sửa.
type Transaction struct {
	Seq       int64 `gorm:"primaryKey;->"` // DB tự sinh, chỉ đọc
	ID        string
	AdminID   string
	Amount    int64 // âm = trừ tiền
	Note      string
	CreatedBy string
	CreatedAt time.Time
}

func (Transaction) TableName() string { return "transactions" }
