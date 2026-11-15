package repository

import (
	"context"
	"errors"

	"github.com/carlosEA28/benefitGuide/internal/catalog/internal/models"
	"github.com/carlosEA28/benefitGuide/internal/shared/exceptions"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RequirementRepository persists the access requirements of a public service.
type RequirementRepository interface {
	CreateRequirement(ctx context.Context, requirement *models.ServiceRequirement) (*models.ServiceRequirement, error)
	FindRequirementByID(ctx context.Context, id uuid.UUID) (*models.ServiceRequirement, error)
	FindByPublicServiceID(ctx context.Context, publicServiceID uuid.UUID) ([]models.ServiceRequirement, error)
	UpdateRequirement(ctx context.Context, requirement *models.ServiceRequirement) (*models.ServiceRequirement, error)
	DeleteRequirement(ctx context.Context, id uuid.UUID) error
}

var _ RequirementRepository = (*repository)(nil)

func (r *repository) CreateRequirement(ctx context.Context, requirement *models.ServiceRequirement) (*models.ServiceRequirement, error) {
	if err := r.db.WithContext(ctx).Create(requirement).Error; err != nil {
		return nil, err
	}
	return requirement, nil
}

func (r *repository) FindRequirementByID(ctx context.Context, id uuid.UUID) (*models.ServiceRequirement, error) {
	var requirement models.ServiceRequirement
	if err := r.db.WithContext(ctx).First(&requirement, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrRequirementNotFound
		}
		return nil, err
	}
	return &requirement, nil
}

func (r *repository) FindByPublicServiceID(ctx context.Context, publicServiceID uuid.UUID) ([]models.ServiceRequirement, error) {
	requirements := make([]models.ServiceRequirement, 0)
	if err := r.db.WithContext(ctx).
		Where("public_service_id = ?", publicServiceID).
		Order("display_order ASC, created_at ASC, id ASC").
		Find(&requirements).Error; err != nil {
		return nil, err
	}
	return requirements, nil
}

func (r *repository) UpdateRequirement(ctx context.Context, requirement *models.ServiceRequirement) (*models.ServiceRequirement, error) {
	if requirement.ID == uuid.Nil {
		return nil, exceptions.ErrRequirementNotFound
	}

	result := r.db.WithContext(ctx).
		Model(requirement).
		Select("description", "type", "condition_description", "display_order").
		Updates(requirement)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.ErrRequirementNotFound
	}
	return requirement, nil
}

func (r *repository) DeleteRequirement(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return exceptions.ErrRequirementNotFound
	}

	result := r.db.WithContext(ctx).Delete(&models.ServiceRequirement{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return exceptions.ErrRequirementNotFound
	}
	return nil
}
