package repository

import (
	"context"
	"errors"

	"github.com/carlosEA28/benefitGuide/internal/catalog/internal/models"
	"github.com/carlosEA28/benefitGuide/internal/shared/exceptions"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PublicServiceRepository reads public services for catalog business rules.
type PublicServiceRepository interface {
	FindPublicServiceByID(ctx context.Context, id uuid.UUID) (*models.PublicService, error)
}

var _ PublicServiceRepository = (*repository)(nil)

func (r *repository) FindPublicServiceByID(ctx context.Context, id uuid.UUID) (*models.PublicService, error) {
	if id == uuid.Nil {
		return nil, exceptions.ErrPublicServiceNotFound
	}

	var publicService models.PublicService
	if err := r.db.WithContext(ctx).First(&publicService, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrPublicServiceNotFound
		}
		return nil, err
	}
	return &publicService, nil
}
