# wallet-service  (:50051)

gRPC `wallet.v1.WalletService` (GetWallet, ListTransactions, AdjustBalance). Authz 2 lớp: interceptor theo RPC (`cmd/grpc/routes/policies.go`) và quyền theo từng ví (`internal/service/wallet/access.go`).

Cần auth-service (JWKS) để verify token: `oauth.jwks_url`. Là chủ của hợp đồng `api/proto/wallet/v1/wallet.proto`.

## Chạy

Folder này là 1 project Go độc lập (module riêng, không cần `go.work`), mở folder lên rồi:

```bash
go run ./cmd/grpc    # hoặc: make run
make test             # go test -race ./...
```

Cấu hình ở `config/config.yml`, ghi đè bằng biến môi trường theo quy tắc `.` → `_`
(ví dụ `REST_SERVER_ADDR=:9999`, `OAUTH_JWKS_URL=...`).

## Docker

```bash
make build            # docker build -t go-admin-wallet/wallet-service:dev-latest .
make up               # chạy riêng service này (gọi service khác qua host.docker.internal)
make ci-docker-build ci-docker-push REGISTRY=<registry>/<group> TAG=v1.0.0
```

## Cấu trúc

```
cmd/grpc/main.go, routes/   # bootstrap + đăng ký route
config/                      # config.go (viper) + config.yml
pkg/authkit                  # verify JWT (JWKS), middleware REST, interceptor gRPC
internal/
    domain/
    internal/repository/{wallet,transaction,memstore}  # Tx khi đổi số dư
    internal/service/wallet
    internal/handler/grpc/v1/wallet
```

## gRPC contract

`api/proto/wallet/v1/wallet.proto` → code sinh ra ở `pkg/gen/wallet/v1`. Đây là bản gốc của hợp đồng gRPC.

```bash
make gen-proto        # buf generate trong container go-tools
# hoặc cài buf + protoc-gen-go + protoc-gen-go-grpc rồi: cd api && buf generate
```
