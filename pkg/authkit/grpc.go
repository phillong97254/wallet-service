package authkit

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor: phía server gRPC. Lấy token từ metadata "authorization",
// verify (authn), rồi tra policy theo tên method (authz).
// Method không có trong policies => từ chối (deny by default).
func UnaryServerInterceptor(v *Verifier, policies map[string]Policy) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		var raw string
		if vals := md.Get("authorization"); len(vals) > 0 {
			raw, _ = strings.CutPrefix(vals[0], "Bearer ")
		}
		if raw == "" {
			return nil, status.Error(codes.Unauthenticated, "missing bearer token")
		}
		claims, err := v.Verify(raw)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token: "+err.Error())
		}
		p, ok := policies[info.FullMethod]
		if !ok || !p.Allows(claims) {
			return nil, status.Errorf(codes.PermissionDenied, "permission denied for %s", info.FullMethod)
		}
		return handler(WithClaims(ctx, claims, raw), req)
	}
}

// ForwardTokenInterceptor: phía client gRPC. Tự gắn access token đang có trong
// context (token của admin gửi tới backoffice) vào metadata khi gọi service khác.
func ForwardTokenInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if raw := RawTokenFrom(ctx); raw != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+raw)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
