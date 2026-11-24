package wallet

import (
	"log/slog"

	walletCache "github.com/phillong97254/wallet-service/internal/cache/wallet"
	"github.com/phillong97254/wallet-service/internal/repository"
)

type impl struct {
	logger *slog.Logger
	repo   repository.Registry
	cache  walletCache.Cache
}

func New(logger *slog.Logger, repo repository.Registry, cache walletCache.Cache) Service {
	return &impl{logger: logger.With("service", "wallet"), repo: repo, cache: cache}
}
