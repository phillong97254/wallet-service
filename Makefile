# Variables
ENV:=dev
GROUP_NAME:=phillong97254
PROJECT_NAME:=wallet-service
REGISTRY:=ghcr.io/phillong97254
TAG:=dev-latest
IMAGE:=$(REGISTRY)/$(PROJECT_NAME)
# Thư viện dùng chung (giống go-library ở các dự án trước)
GO_LIB:=github.com/phillong97254/go-lib
GO_LIB_VERSION:=latest
GO_LIB_LOCAL:=../go-lib
GOPRIVATE:=github.com/phillong97254/*
# Shorten cmd
DOCKER_BIN := docker
COMPOSE_BIN := ENV=$(ENV) GROUP_NAME=$(GROUP_NAME) PROJECT_NAME=$(PROJECT_NAME) REGISTRY=$(REGISTRY) TAG=$(TAG) docker compose
COMPOSE_GO_TOOL_RUN := $(COMPOSE_BIN) -f docker-compose.tools.yml run --rm go-tools

.PHONY: netrc run mod mod-local mod-remote test build-bin up down logs build docker-build docker-login docker-push \
	ci-docker-build ci-docker-login ci-docker-push ci-docker-clean gen-proto

# ---------- Go (local) ----------
run:
	go run ./cmd/grpc

# Credential cho repo private (go-lib): thêm dòng "machine github.com login phillong97254 password <personal-access-token>"
netrc:
	touch ~/.netrc
	chmod 600 ~/.netrc

# Cập nhật thư viện: go get -u + tidy + vendor (giống các dự án trước).
# go-lib: nếu go.mod đang replace sang $(GO_LIB_LOCAL) thì dùng code local, ngược lại
# go get theo GO_LIB_VERSION (make mod GO_LIB_VERSION=v0.2.0 để chốt version).
mod:
	@if grep -q '$(GO_LIB) =>' go.mod; then \
		echo "go-lib: dùng bản local $(GO_LIB_LOCAL)"; \
		go list -m -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' all | grep -v '^$(GO_LIB)$$' | xargs -r go get -u; \
	else \
		GOPRIVATE=$(GOPRIVATE) go get -u ./... && GOPRIVATE=$(GOPRIVATE) go get $(GO_LIB)@$(GO_LIB_VERSION); \
	fi
	go mod tidy
	go mod vendor

# Dùng go-lib local (chưa push / đang sửa go-lib cùng lúc)
mod-local:
	go mod edit -replace=$(GO_LIB)=$(GO_LIB_LOCAL)
	@$(MAKE) mod

# Dùng go-lib từ remote theo version (sau khi go-lib đã push + tag)
mod-remote:
	go mod edit -dropreplace=$(GO_LIB)
	@$(MAKE) mod

test:
	go test -race ./...

build-bin:
	go build -o bin/$(PROJECT_NAME) ./cmd/grpc

# ---------- Docker compose (chỉ service này) ----------
up:
	@$(COMPOSE_BIN) up -d

down:
	@$(COMPOSE_BIN) down

logs:
	@$(COMPOSE_BIN) logs -f

# ---------- Docker image ----------
build:
	-@$(DOCKER_BIN) image rm -f $(IMAGE):dev-latest
	@$(DOCKER_BIN) build --rm --force-rm -t $(IMAGE):dev-latest .
	-@$(DOCKER_BIN) image prune -f --filter label=stage=Builder

docker-build:
	@$(DOCKER_BIN) build --rm --force-rm -t $(IMAGE):dev-latest .
	-@$(DOCKER_BIN) image prune -f --filter label=stage=Builder

# Đăng nhập registry tương tác (không lưu mật khẩu trong Makefile)
docker-login:
	@$(DOCKER_BIN) login $(firstword $(subst /, ,$(REGISTRY)))

docker-push: docker-build docker-login
	@$(DOCKER_BIN) push $(IMAGE):dev-latest

ci-docker-build:
	@$(DOCKER_BIN) build --rm --force-rm -t $(IMAGE):$(TAG) .

ci-docker-login:
	@echo "$(DOCKER_TOKEN)" | $(DOCKER_BIN) login $(firstword $(subst /, ,$(REGISTRY))) -u $(DOCKER_USER) --password-stdin

ci-docker-push:
	@$(DOCKER_BIN) push $(IMAGE):$(TAG)

ci-docker-clean:
	@$(DOCKER_BIN) image rm -f $(IMAGE):$(TAG)
	-@$(DOCKER_BIN) image prune -f

# ---------- Codegen ----------
# Sinh lại code gRPC từ api/proto vào pkg/gen bằng image go-tools (không cần cài buf trên máy)
gen-proto:
	@$(COMPOSE_GO_TOOL_RUN) sh -c 'cd api && buf generate'
