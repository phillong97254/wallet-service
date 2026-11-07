// wallet-service: gRPC server giữ thông tin ví admin.
// Authz 2 lớp:
//  1. Interceptor: scope/role theo từng RPC (cmd/grpc/routes/policies.go).
//  2. Service: quyền trên tài nguyên cụ thể (admin chỉ xem ví của mình).
package main

import (
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	auth "github.com/phillong97254/go-lib/authkit"
	"github.com/phillong97254/go-lib/database"
	"github.com/phillong97254/wallet-service/cmd/grpc/routes"
	"github.com/phillong97254/wallet-service/config"
	v1 "github.com/phillong97254/wallet-service/internal/handler/grpc/v1"
	"github.com/phillong97254/wallet-service/internal/repository"
	"github.com/phillong97254/wallet-service/internal/service/wallet"
	"github.com/phillong97254/wallet-service/migration"
)

func main() {
	// Load config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With("server", cfg.GRPCServer.Name, "env", cfg.API.Env)
	slog.SetDefault(logger)

	// Authorize: verify JWT bằng JWKS của auth-service
	verifier := auth.NewVerifier(cfg.OAuth.JWKSURL, cfg.OAuth.Issuer, cfg.OAuth.Audience)

	// Init Database (+ migration)
	db, err := database.New(database.Config{
		URL:                   cfg.DB.URL,
		MaxOpenConnections:    cfg.DB.MaxOpenConnections,
		MaxIdleConnections:    cfg.DB.MaxIdleConnections,
		ConnectionMaxLifetime: cfg.DB.ConnectionMaxLifetime,
		ConnectionMaxIdleTime: cfg.DB.ConnectionMaxIdleTime,
		LogSQL:                cfg.DB.LogSQL,
	})
	if err != nil {
		log.Fatal(err)
	}
	if cfg.DB.AutoMigrate {
		if err := database.Migrate(cfg.DB.URL, migration.FS, migration.Dir); err != nil {
			log.Fatal(err)
		}
	}

	// Repository -> Service -> Handler
	repo := repository.New(db)
	walletService := wallet.New(logger, repo)
	handler := v1.New(logger, walletService)

	// Register Routes
	srv := grpc.NewServer(grpc.UnaryInterceptor(auth.UnaryServerInterceptor(verifier, routes.Policies)))
	routes.V1Routes(srv, handler)

	// Start Server
	lis, err := net.Listen("tcp", cfg.GRPCServer.Addr)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		logger.Info("wallet-service (gRPC) listening", "addr", cfg.GRPCServer.Addr)
		if err := srv.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	srv.GracefulStop()
	logger.Info("wallet-service stopped")
}
