package wallet

import (
	"context"

	"github.com/phillong97254/wallet-service/internal/domain"
)

// nop là cache tắt: Get luôn miss, Set/Delete không làm gì.
type nop struct{}

func (nop) Get(context.Context, string) (domain.Wallet, error) { return domain.Wallet{}, ErrMiss }
func (nop) Set(context.Context, domain.Wallet) error           { return nil }
func (nop) Delete(context.Context, string) error               { return nil }
