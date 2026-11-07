package transaction

import "github.com/phillong97254/wallet-service/internal/repository/memstore"

type impl struct {
	store *memstore.Store
}

func New(store *memstore.Store) Repository {
	return &impl{store: store}
}
