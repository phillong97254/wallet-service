// Package authkit gom mọi thứ liên quan OAuth2/JWT dùng chung cho các service:
// verify token qua JWKS, HTTP middleware, gRPC interceptor, kiểm tra scope/role.
package authkit

import (
	"context"
	"slices"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Claims là payload của access token do auth-service (hoặc Keycloak) phát hành.
type Claims struct {
	Scope    string   `json:"scope"`               // "wallet:read users:read" (chuẩn OAuth2: cách nhau bằng dấu cách)
	Roles    []string `json:"roles,omitempty"`     // role của người dùng: admin, superadmin
	ClientID string   `json:"client_id,omitempty"` // client nào xin token
	jwt.RegisteredClaims
}

func (c *Claims) HasScope(s string) bool { return slices.Contains(strings.Fields(c.Scope), s) }

func (c *Claims) HasAnyRole(roles ...string) bool {
	for _, r := range roles {
		if slices.Contains(c.Roles, r) {
			return true
		}
	}
	return false
}

// Policy mô tả quyền cần có cho 1 endpoint / 1 RPC.
type Policy struct {
	Scopes []string // phải có TẤT CẢ các scope này
	Roles  []string // phải có ÍT NHẤT 1 role (để trống = không yêu cầu role)
}

func (p Policy) Allows(c *Claims) bool {
	for _, s := range p.Scopes {
		if !c.HasScope(s) {
			return false
		}
	}
	return len(p.Roles) == 0 || c.HasAnyRole(p.Roles...)
}

type ctxKey int

const (
	claimsKey ctxKey = iota
	tokenKey
)

func WithClaims(ctx context.Context, c *Claims, raw string) context.Context {
	ctx = context.WithValue(ctx, claimsKey, c)
	return context.WithValue(ctx, tokenKey, raw)
}

func ClaimsFrom(ctx context.Context) (*Claims, bool) {
	c, ok := ctx.Value(claimsKey).(*Claims)
	return c, ok
}

// RawTokenFrom trả về access token gốc, dùng để forward sang service khác.
func RawTokenFrom(ctx context.Context) string {
	s, _ := ctx.Value(tokenKey).(string)
	return s
}
