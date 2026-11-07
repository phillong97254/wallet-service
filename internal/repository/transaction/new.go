package transaction

import "gorm.io/gorm"

type impl struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &impl{db: db}
}
