// Package database mở kết nối Postgres (gorm) và chạy migration (golang-migrate),
// dùng chung cho các service (giống database.New ở go-library của các dự án trước).
package database

import (
	"errors"
	"fmt"
	"io/fs"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5" // driver "pgx5://"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	URL                   string // postgres://user:pass@host:5432/db?sslmode=disable
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime int // phút
	ConnectionMaxIdleTime int // phút
	LogSQL                bool
}

// New mở *gorm.DB và kiểm tra kết nối (fail fast). Thử lại vài lần vì khi chạy docker compose
// Postgres có thể chưa sẵn sàng.
func New(cfg Config) (*gorm.DB, error) {
	if cfg.URL == "" {
		return nil, errors.New("database: missing URL")
	}
	logLevel := logger.Warn
	if cfg.LogSQL {
		logLevel = logger.Info
	}
	db, err := gorm.Open(postgres.Open(cfg.URL), &gorm.Config{
		Logger:         logger.Default.LogMode(logLevel),
		TranslateError: true, // trả gorm.ErrDuplicatedKey, ...
		NowFunc:        func() time.Time { return time.Now().UTC() },
	})
	if err != nil {
		return nil, fmt.Errorf("database: open: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}
	if cfg.MaxOpenConnections > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConnections)
	}
	if cfg.MaxIdleConnections > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	}
	if cfg.ConnectionMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnectionMaxLifetime) * time.Minute)
	}
	if cfg.ConnectionMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnectionMaxIdleTime) * time.Minute)
	}

	for attempt := 1; ; attempt++ {
		if err = sqlDB.Ping(); err == nil {
			return db, nil
		}
		if attempt == 10 {
			return nil, fmt.Errorf("database: ping: %w", err)
		}
		time.Sleep(time.Second)
	}
}

// Migrate chạy mọi migration "up" trong thư mục dir của migrations (embed.FS).
// File đặt tên theo golang-migrate: 000000001_init.up.sql / 000000001_init.down.sql.
func Migrate(url string, migrations fs.FS, dir string) error {
	src, err := iofs.New(migrations, dir)
	if err != nil {
		return fmt.Errorf("migrate: source: %w", err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", src, pgx5URL(url))
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	defer m.Close()
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate: up: %w", err)
	}
	return nil
}

// pgx5URL đổi postgres:// thành pgx5:// (scheme của driver pgx v5 trong golang-migrate).
func pgx5URL(url string) string {
	for _, p := range []string{"postgres://", "postgresql://"} {
		if len(url) > len(p) && url[:len(p)] == p {
			return "pgx5://" + url[len(p):]
		}
	}
	return url
}
