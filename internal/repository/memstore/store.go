// Package memstore là "database" in-memory dùng chung cho các repository.
// Đổi sang Postgres thì thay package này bằng *gorm.DB / *sql.DB.
package memstore

import (
	"sync"
	"time"

	"github.com/phil/wallet-service/internal/domain"
)

type Store struct {
	Mu           sync.RWMutex
	Wallets      map[string]domain.Wallet
	Transactions map[string][]domain.Transaction // theo AdminID, cũ -> mới
}

func New() *Store {
	now := time.Now().UTC()
	return &Store{
		Wallets: map[string]domain.Wallet{
			"alice": {AdminID: "alice", Currency: "VND", Balance: 5_000_000, UpdatedAt: now},
			"bob":   {AdminID: "bob", Currency: "VND", Balance: 12_000_000, UpdatedAt: now},
		},
		Transactions: map[string][]domain.Transaction{
			"alice": {{ID: "tx-1", AdminID: "alice", Amount: 5_000_000, Note: "initial funding", CreatedBy: "system", CreatedAt: now}},
			"bob":   {{ID: "tx-1", AdminID: "bob", Amount: 12_000_000, Note: "initial funding", CreatedBy: "system", CreatedAt: now}},
		},
	}
}
