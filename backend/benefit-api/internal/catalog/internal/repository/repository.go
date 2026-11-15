package repository

import (
	"context"

	"gorm.io/gorm"
)

// Requirement defines the access requirements for a service/repository.
// Implementations can check specific database conditions, tables, or configurations
// that the service needs to function properly.
type Requirement interface {
	// Ready checks if the underlying connection/data is healthy.
	Ready(ctx context.Context) error
	// Check validates that the specific requirement is satisfied.
	// Returns nil if the requirement is met, or an error describing the failure.
	Check(ctx context.Context) error
}

// Repository base struct with requirement support.
// Embedding Requirement allows Repository to satisfy the Requirement interface.
type Repository struct {
	db *gorm.DB
}

// New creates a new Repository instance.
func New(db *gorm.DB) *Repository { return &Repository{db: db} }

// Ready checks if the database connection is healthy.
func (r *Repository) Ready(ctx context.Context) error {
	return r.db.WithContext(ctx).Exec("SELECT 1").Error
}

// Check implements the Requirement interface and validates basic repository conditions.
// By default, it checks database connectivity. Subtypes can override this to add
// specific requirements (e.g., specific tables, columns, or configurations).
func (r *Repository) Check(ctx context.Context) error {
	if err := r.db.WithContext(ctx).Exec("SELECT 1").Error; err != nil {
		return err
	}
	return nil
}
