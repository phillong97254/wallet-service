package wallet

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"sync"
	"testing"

	"github.com/golang-jwt/jwt/v5"

	auth "github.com/phillong97254/go-lib/authkit"
	"github.com/phillong97254/wallet-service/internal/repository"
	"github.com/phillong97254/wallet-service/internal/repository/memstore"
)

func newTestService() Service {
	return New(slog.New(slog.NewTextHandler(io.Discard, nil)), repository.New(memstore.New()))
}

func ctxAs(sub string, roles ...string) context.Context {
	c := &auth.Claims{Roles: roles, RegisteredClaims: jwt.RegisteredClaims{Subject: sub}}
	return auth.WithClaims(context.Background(), c, "test-token")
}

func TestAdjustBalance(t *testing.T) {
	tests := []struct {
		name    string
		adminID string
		amount  int64
		wantErr error
		wantBal int64
	}{
		{name: "cộng tiền", adminID: "alice", amount: 250_000, wantBal: 5_250_000},
		{name: "trừ tiền", adminID: "alice", amount: -1_000_000, wantBal: 4_000_000},
		{name: "không đủ số dư", adminID: "alice", amount: -6_000_000, wantErr: ErrInsufficientBalance},
		{name: "amount = 0", adminID: "alice", amount: 0, wantErr: ErrInvalidAmount},
		{name: "ví không tồn tại", adminID: "zed", amount: 1, wantErr: ErrNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService()
			out, err := svc.AdjustBalance(ctxAs("bob", "admin", "superadmin"), AdjustBalanceInput{AdminID: tt.adminID, Amount: tt.amount})
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && out.Balance != tt.wantBal {
				t.Fatalf("balance = %d, want %d", out.Balance, tt.wantBal)
			}
		})
	}
}

func TestAdjustBalanceRecordsAudit(t *testing.T) {
	svc := newTestService()
	ctx := ctxAs("bob", "admin", "superadmin")
	if _, err := svc.AdjustBalance(ctx, AdjustBalanceInput{AdminID: "alice", Amount: 100, Note: "bonus"}); err != nil {
		t.Fatal(err)
	}
	txs, err := svc.ListTransactions(ctx, ListTransactionsInput{AdminID: "alice"})
	if err != nil {
		t.Fatal(err)
	}
	last := txs.Items[len(txs.Items)-1]
	if last.ID != "tx-2" || last.CreatedBy != "bob" || last.Note != "bonus" {
		t.Fatalf("last tx = %+v", last)
	}
}

// Nhiều request trừ tiền song song không được làm số dư âm (Tx tuần tự hoá đọc-kiểm tra-ghi).
func TestAdjustBalanceConcurrent(t *testing.T) {
	svc := newTestService()
	ctx := ctxAs("bob", "admin", "superadmin")
	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = svc.AdjustBalance(ctx, AdjustBalanceInput{AdminID: "alice", Amount: -100_000})
		}()
	}
	wg.Wait()
	w, err := svc.GetWallet(ctx, GetWalletInput{AdminID: "alice"})
	if err != nil {
		t.Fatal(err)
	}
	if w.Balance != 0 {
		t.Fatalf("balance = %d, want 0 (5_000_000 / 100_000 = 50 lần thành công)", w.Balance)
	}
}

func TestGetWalletAccess(t *testing.T) {
	svc := newTestService()
	if _, err := svc.GetWallet(ctxAs("alice", "admin"), GetWalletInput{AdminID: "alice"}); err != nil {
		t.Fatalf("alice xem ví mình: %v", err)
	}
	if _, err := svc.GetWallet(ctxAs("alice", "admin"), GetWalletInput{AdminID: "bob"}); !errors.Is(err, ErrForbidden) {
		t.Fatalf("alice xem ví bob: err = %v, want ErrForbidden", err)
	}
	if _, err := svc.GetWallet(ctxAs("bob", "admin", "superadmin"), GetWalletInput{AdminID: "alice"}); err != nil {
		t.Fatalf("superadmin xem ví alice: %v", err)
	}
}
