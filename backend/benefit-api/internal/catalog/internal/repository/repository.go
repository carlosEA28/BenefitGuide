package repository

import "gorm.io/gorm"

// repository persists catalog entities. It is unexported on purpose: callers
// depend on the per-aggregate interfaces declared next to each entity instead.
type repository struct {
	db *gorm.DB
}

// New creates a new repository.
func New(db *gorm.DB) *repository { return &repository{db: db} }
