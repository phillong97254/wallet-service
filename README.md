# wallet-service  (:50051)

gRPC `wallet.v1.WalletService` (GetWallet, ListTransactions, AdjustBalance). Authz 2 lớp: interceptor theo RPC (`cmd/grpc/routes/policies.go`) và quyền theo từng ví (`internal/service/wallet/access.go`).

Cần auth-service (JWKS) để verify token: `oauth.jwks_url`. Là chủ của hợp đồng `api/proto/wallet/v1/wallet.proto`.

## Chạy

Folder này là 1 project Go độc lập (repo git, module riêng), mở folder lên rồi:

```bash
make env              # lần đầu: tạo .env từ .env.example
go run ./cmd/grpc    # hoặc: make run
make test             # go test -race ./...
```

Cấu hình ở `config/config.yml`, ghi đè bằng biến môi trường theo quy tắc `.` → `_`
(ví dụ `REST_SERVER_ADDR=:9999`, `OAUTH_JWKS_URL=...`).

## Secret và cấu hình

- `config/config.yml` **được commit**: chỉ chứa giá trị mặc định, không chứa secret.
- Secret đặt trong `.env` (đã có trong `.gitignore` và `.dockerignore`, không commit, không vào image).
  `make env` tạo `.env` từ `.env.example`, sau đó điền giá trị. `make run` tự nạp `.env`.
- Thiếu secret thì service dừng ngay lúc khởi động và báo tên biến cần đặt.
- Trên server/CI: truyền secret bằng biến môi trường (GitHub Actions secrets, Docker/K8s secret...), không dùng `.env` commit lên.

## Thư viện (make mod)

Code dùng chung (`authkit`, `apperr`) nằm ở repo **go-lib** (`github.com/phillong97254/go-lib`), giống `go-library`
ở các dự án trước. Thư viện được vendor vào `vendor/` nên build/Docker không cần mạng hay repo khác.

```bash
make mod                          # go get -u thư viện + go mod tidy + go mod vendor
make mod GO_LIB_VERSION=v0.2.0    # chốt go-lib ở 1 version (khi dùng remote)
make mod-local                    # dùng ../go-lib (replace) khi đang sửa go-lib cùng lúc
make mod-remote                   # bỏ replace, lấy go-lib theo version từ GitHub (mặc định hiện tại)
```

## Docker

```bash
make build            # docker build -t ghcr.io/phillong97254/wallet-service:dev-latest .
make up               # chạy riêng service này (gọi service khác qua host.docker.internal)
make ci-docker-build ci-docker-push TAG=v1.0.0   # push lên ghcr.io/phillong97254 (GitHub Container Registry)
```

## Cấu trúc

```
cmd/grpc/main.go, routes/   # bootstrap + đăng ký route
config/                      # config.go (viper) + config.yml
vendor/                      # thư viện (kể cả go-lib), tạo bằng make mod
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
