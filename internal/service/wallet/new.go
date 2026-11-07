package wallet

import (
	"log/slog"

	"github.com/phil/wallet-service/internal/repository"
)

type impl struct {
	logger *slog.Logger
	repo   repository.Registry
}

func New(logger *slog.Logger, repo repository.Registry) Service {
	return &impl{logger: logger.With("service", "wallet"), repo: repo}
}
