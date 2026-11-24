package wallet

import (
	"context"
	"errors"
	"sync"
	"testing"

	walletCache "github.com/phillong97254/wallet-service/internal/cache/wallet"
	"github.com/phillong97254/wallet-service/internal/domain"
)

// fakeCache là cache trong bộ nhớ, đếm số lần gọi để test biết có đọc/xoá cache hay không.
type fakeCache struct {
	mu      sync.Mutex
	data    map[string]domain.Wallet
	getErr  error
	setErr  error
	delErr  error
	hits    int
	deletes int
}

func newFakeCache() *fakeCache {
	return &fakeCache{data: map[string]domain.Wallet{}}
}

func (f *fakeCache) Get(_ context.Context, adminID string) (domain.Wallet, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.getErr != nil {
		return domain.Wallet{}, f.getErr
	}
	w, ok := f.data[adminID]
	if !ok {
		return domain.Wallet{}, walletCache.ErrMiss
	}
	f.hits++
	return w, nil
}

func (f *fakeCache) Set(_ context.Context, w domain.Wallet) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.setErr != nil {
		return f.setErr
	}
	f.data[w.AdminID] = w
	return nil
}

func (f *fakeCache) Delete(_ context.Context, adminID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.deletes++
	if f.delErr != nil {
		return f.delErr
	}
	delete(f.data, adminID)
	return nil
}

func TestGetWallet_Cache(t *testing.T) {
	ctx := ctxAs("alice", "admin")

	t.Run("miss rồi nạp cache, lần sau hit", func(t *testing.T) {
		cache := newFakeCache()
		svc := newTestServiceWithCache(t, cache)

		first, err := svc.GetWallet(ctx, GetWalletInput{AdminID: "alice"})
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := cache.data["alice"]; !ok {
			t.Fatal("ví chưa được nạp vào cache sau lần đọc đầu")
		}
		second, err := svc.GetWallet(ctx, GetWalletInput{AdminID: "alice"})
		if err != nil {
			t.Fatal(err)
		}
		if cache.hits != 1 || second.Balance != first.Balance {
			t.Fatalf("hits = %d, balance %d vs %d", cache.hits, second.Balance, first.Balance)
		}
	})

	t.Run("AdjustBalance xoá cache, lần đọc sau thấy số dư mới", func(t *testing.T) {
		cache := newFakeCache()
		svc := newTestServiceWithCache(t, cache)

		before, err := svc.GetWallet(ctx, GetWalletInput{AdminID: "alice"})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := svc.AdjustBalance(ctxAs("bob", "admin", "superadmin"), AdjustBalanceInput{AdminID: "alice", Amount: 1_000}); err != nil {
			t.Fatal(err)
		}
		if cache.deletes != 1 {
			t.Fatalf("deletes = %d, want 1", cache.deletes)
		}
		after, err := svc.GetWallet(ctx, GetWalletInput{AdminID: "alice"})
		if err != nil {
			t.Fatal(err)
		}
		if after.Balance != before.Balance+1_000 {
			t.Fatalf("balance = %d, want %d", after.Balance, before.Balance+1_000)
		}
	})

	t.Run("AdjustBalance lỗi nghiệp vụ không xoá cache", func(t *testing.T) {
		cache := newFakeCache()
		svc := newTestServiceWithCache(t, cache)

		_, err := svc.AdjustBalance(ctxAs("bob", "admin", "superadmin"), AdjustBalanceInput{AdminID: "alice", Amount: -999_999_999})
		if !errors.Is(err, ErrInsufficientBalance) {
			t.Fatalf("err = %v", err)
		}
		if cache.deletes != 0 {
			t.Fatalf("deletes = %d, want 0", cache.deletes)
		}
	})

	t.Run("cache lỗi thì vẫn đọc được từ DB", func(t *testing.T) {
		cache := newFakeCache()
		cache.getErr = errors.New("redis down")
		cache.setErr = errors.New("redis down")
		svc := newTestServiceWithCache(t, cache)

		out, err := svc.GetWallet(ctx, GetWalletInput{AdminID: "alice"})
		if err != nil || out.AdminID != "alice" {
			t.Fatalf("out = %+v, err = %v", out, err)
		}
	})

	t.Run("xoá cache lỗi thì AdjustBalance vẫn thành công", func(t *testing.T) {
		cache := newFakeCache()
		cache.delErr = errors.New("redis down")
		svc := newTestServiceWithCache(t, cache)

		if _, err := svc.AdjustBalance(ctxAs("bob", "admin", "superadmin"), AdjustBalanceInput{AdminID: "alice", Amount: 1}); err != nil {
			t.Fatalf("err = %v", err)
		}
	})

	t.Run("không có quyền thì không đọc cache", func(t *testing.T) {
		cache := newFakeCache()
		cache.data["bob"] = domain.Wallet{AdminID: "bob", Balance: 1}
		svc := newTestServiceWithCache(t, cache)

		_, err := svc.GetWallet(ctx, GetWalletInput{AdminID: "bob"})
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("err = %v, want ErrForbidden", err)
		}
		if cache.hits != 0 {
			t.Fatalf("hits = %d, want 0", cache.hits)
		}
	})

	t.Run("ví không tồn tại không bị cache", func(t *testing.T) {
		cache := newFakeCache()
		svc := newTestServiceWithCache(t, cache)

		_, err := svc.GetWallet(ctxAs("bob", "admin", "superadmin"), GetWalletInput{AdminID: "zed"})
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("err = %v, want ErrNotFound", err)
		}
		if len(cache.data) != 0 {
			t.Fatalf("cache có %d key, want 0", len(cache.data))
		}
	})
}
