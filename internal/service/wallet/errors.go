package wallet

import "errors"

// Lỗi nghiệp vụ, handler gRPC đổi sang status code tương ứng.
var (
	ErrUnauthenticated     = errors.New("unauthenticated")
	ErrForbidden           = errors.New("you can only access your own wallet")
	ErrNotFound            = errors.New("wallet not found")
	ErrInvalidAmount       = errors.New("amount must not be zero")
	ErrInsufficientBalance = errors.New("insufficient balance")
)
