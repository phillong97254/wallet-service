package wallet

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type impl struct {
	client *redis.Client
	ttl    time.Duration
}

// New trả cache Redis; addr rỗng thì trả cache tắt (luôn miss) để service chạy được khi chưa có Redis.
// Timeout ngắn để Redis chậm/chết không kéo theo request đọc ví.
func New(addr, password string, ttl time.Duration) Cache {
	if addr == "" {
		return nop{}
	}
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DialTimeout:  time.Second,
		ReadTimeout:  500 * time.Millisecond,
		WriteTimeout: 500 * time.Millisecond,
		MaxRetries:   1,
	})
	return &impl{client: client, ttl: ttl}
}

func key(adminID string) string {
	return "wallet:" + adminID
}
