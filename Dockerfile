# Build từ chính thư mục service:  docker build -t go-admin-wallet/wallet-service:dev-latest .
FROM golang:1.26.1 AS builder
LABEL stage=Builder

ENV GROUP_NAME=phil
ENV GOPRIVATE=github.com/phil/*
ENV PROJECT_NAME=wallet-service

WORKDIR /app
COPY ./vendor ./vendor
COPY ./go.mod .
COPY ./go.sum .
COPY ./config ./config
COPY ./pkg ./pkg
COPY ./internal ./internal
COPY ./cmd ./cmd

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags="-s -w" -installsuffix cgo -o main ./cmd/grpc


FROM alpine:latest AS final
LABEL stage=Final
RUN apk --no-cache add ca-certificates tzdata \
    && cp /usr/share/zoneinfo/UTC /etc/localtime \
    && echo "UTC" > /etc/timezone \
    && apk del tzdata
WORKDIR /root

# Copy the Pre-built binary file + config from the previous stage
COPY --from=builder /app/main .
COPY --from=builder /app/config/config.yml ./config/config.yml

EXPOSE 50051

#Command to run the executable
CMD ["./main"]
