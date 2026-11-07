package domain

import "time"

// Transaction là một lần thay đổi số dư. CreatedBy là audit: ai đã sửa.
type Transaction struct {
	ID        string
	AdminID   string
	Amount    int64 // âm = trừ tiền
	Note      string
	CreatedBy string
	CreatedAt time.Time
}
