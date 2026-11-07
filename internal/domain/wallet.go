package domain

import "time"

// Wallet là ví của một admin. Balance tính theo đơn vị nhỏ nhất (VND: đồng).
type Wallet struct {
	AdminID   string
	Currency  string
	Balance   int64
	UpdatedAt time.Time
}
