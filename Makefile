# Variables
ENV:=dev
GROUP_NAME:=phil
PROJECT_NAME:=wallet-service
REGISTRY:=go-admin-wallet
TAG:=dev-latest
IMAGE:=$(REGISTRY)/$(PROJECT_NAME)
# Shorten cmd
DOCKER_BIN := docker
COMPOSE_BIN := ENV=$(ENV) GROUP_NAME=$(GROUP_NAME) PROJECT_NAME=$(PROJECT_NAME) REGISTRY=$(REGISTRY) TAG=$(TAG) docker compose
COMPOSE_GO_TOOL_RUN := $(COMPOSE_BIN) -f docker-compose.tools.yml run --rm go-tools

.PHONY: run mod test build-bin up down logs build docker-build docker-login docker-push \
	ci-docker-build ci-docker-login ci-docker-push ci-docker-clean gen-proto

# ---------- Go (local) ----------
run:
	go run ./cmd/grpc

mod:
	go mod tidy

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
