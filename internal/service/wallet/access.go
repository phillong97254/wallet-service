package wallet

import (
	"context"

	auth "github.com/phil/wallet-service/pkg/authkit"
)

// canAccess: quyền trên tài nguyên cụ thể. Superadmin xem được mọi ví,
// admin thường chỉ xem ví của chính mình. (Scope/role theo RPC đã check ở interceptor.)
func canAccess(ctx context.Context, adminID string) error {
	c, ok := auth.ClaimsFrom(ctx)
	if !ok {
		return ErrUnauthenticated
	}
	if c.Subject == adminID || c.HasAnyRole("superadmin") {
		return nil
	}
	return ErrForbidden
}
