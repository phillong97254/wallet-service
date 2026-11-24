package wallet

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/phillong97254/wallet-service/internal/domain"
)

// Chạy với Redis thật: đặt REDIS_TEST_ADDR (+ REDIS_TEST_PASSWORD), dùng Redis riêng cho test.
func newRedisCache(t *testing.T, ttl time.Duration) Cache {
	t.Helper()
	addr := os.Getenv("REDIS_TEST_ADDR")
	if addr == "" {
		t.Skip("REDIS_TEST_ADDR chưa đặt, bỏ qua test cần Redis")
	}
	return New(addr, os.Getenv("REDIS_TEST_PASSWORD"), ttl)
}

func TestCache_SetGetDelete(t *testing.T) {
	ctx := context.Background()
	c := newRedisCache(t, time.Minute)
	w := domain.Wallet{AdminID: "cache-test", Currency: "VND", Balance: 123, UpdatedAt: time.Now().UTC().Truncate(time.Second)}

	if _, err := c.Get(ctx, w.AdminID); !errors.Is(err, ErrMiss) {
		t.Fatalf("Get trước Set: err = %v, want ErrMiss", err)
	}
	if err := c.Set(ctx, w); err != nil {
		t.Fatal(err)
	}
	got, err := c.Get(ctx, w.AdminID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Balance != w.Balance || got.Currency != w.Currency || !got.UpdatedAt.Equal(w.UpdatedAt) {
		t.Fatalf("got %+v, want %+v", got, w)
	}
	if err := c.Delete(ctx, w.AdminID); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Get(ctx, w.AdminID); !errors.Is(err, ErrMiss) {
		t.Fatalf("Get sau Delete: err = %v, want ErrMiss", err)
	}
}

func TestCache_Expires(t *testing.T) {
	ctx := context.Background()
	c := newRedisCache(t, time.Second)
	if err := c.Set(ctx, domain.Wallet{AdminID: "cache-ttl", Balance: 1}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(1500 * time.Millisecond)
	if _, err := c.Get(ctx, "cache-ttl"); !errors.Is(err, ErrMiss) {
		t.Fatalf("err = %v, want ErrMiss sau TTL", err)
	}
}

func TestCache_Disabled(t *testing.T) {
	ctx := context.Background()
	c := New("", "", time.Minute)
	if err := c.Set(ctx, domain.Wallet{AdminID: "x"}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Get(ctx, "x"); !errors.Is(err, ErrMiss) {
		t.Fatalf("err = %v, want ErrMiss khi tắt cache", err)
	}
}

func TestCache_RedisDown(t *testing.T) {
	c := New("127.0.0.1:1", "", time.Minute) // cổng không có ai nghe
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := c.Get(ctx, "x")
	if err == nil || errors.Is(err, ErrMiss) {
		t.Fatalf("err = %v, want lỗi kết nối (không phải ErrMiss)", err)
	}
}
