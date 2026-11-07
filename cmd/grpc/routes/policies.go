package routes

import (
	auth "github.com/phil/go-lib/authkit"
	walletv1 "github.com/phil/wallet-service/pkg/gen/wallet/v1"
)

// Policies: authz lớp 1 (coarse-grained), map tên RPC -> scope + role.
// RPC không có trong bảng thì bị từ chối (deny by default).
// Lớp 2 (quyền trên từng ví) nằm ở service/wallet/access.go.
var Policies = map[string]auth.Policy{
	walletv1.WalletService_GetWallet_FullMethodName:        {Scopes: []string{"wallet:read"}, Roles: []string{"admin"}},
	walletv1.WalletService_ListTransactions_FullMethodName: {Scopes: []string{"wallet:read"}, Roles: []string{"admin"}},
	walletv1.WalletService_AdjustBalance_FullMethodName:    {Scopes: []string{"wallet:write"}, Roles: []string{"superadmin"}},
}
