package repository

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func New(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Ready(ctx context.Context) error {
	return r.db.WithContext(ctx).Exec("SELECT 1").Error
}
