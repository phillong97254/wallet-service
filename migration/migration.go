// Package migration nhúng các file SQL (golang-migrate) vào binary để service tự migrate khi khởi động
// (db.auto_migrate). Chạy tay: make migrate.
package migration

import "embed"

//go:embed scripts/*.sql
var FS embed.FS

const Dir = "scripts"
