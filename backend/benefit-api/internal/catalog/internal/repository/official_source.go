package repository

import (
	"context"
	"errors"

	"github.com/carlosEA28/benefitGuide/internal/catalog/internal/models"
	"github.com/carlosEA28/benefitGuide/internal/shared/exceptions"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OfficialSourceRepository persists the official references of public services.
type OfficialSourceRepository interface {
	Create(ctx context.Context, source *models.OfficialSource) (*models.OfficialSource, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.OfficialSource, error)
	FindByPublicServiceID(ctx context.Context, publicServiceID uuid.UUID) ([]models.OfficialSource, error)
	Update(ctx context.Context, source *models.OfficialSource) (*models.OfficialSource, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

var _ OfficialSourceRepository = (*officialSourceRepository)(nil)

type officialSourceRepository struct {
	db *gorm.DB
}

// NewOfficialSourceRepository creates a repository for official sources.
func NewOfficialSourceRepository(db *gorm.DB) OfficialSourceRepository {
	return &officialSourceRepository{db: db}
}

func (r *officialSourceRepository) Create(ctx context.Context, source *models.OfficialSource) (*models.OfficialSource, error) {
	if err := r.db.WithContext(ctx).Create(source).Error; err != nil {
		return nil, err
	}
	return source, nil
}

func (r *officialSourceRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.OfficialSource, error) {
	if id == uuid.Nil {
		return nil, exceptions.ErrOfficialSourceNotFound
	}

	var source models.OfficialSource
	if err := r.db.WithContext(ctx).First(&source, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrOfficialSourceNotFound
		}
		return nil, err
	}
	return &source, nil
}

func (r *officialSourceRepository) FindByPublicServiceID(ctx context.Context, publicServiceID uuid.UUID) ([]models.OfficialSource, error) {
	sources := make([]models.OfficialSource, 0)
	if err := r.db.WithContext(ctx).
		Where("public_service_id = ?", publicServiceID).
		Order("created_at ASC, id ASC").
		Find(&sources).Error; err != nil {
		return nil, err
	}
	return sources, nil
}

func (r *officialSourceRepository) Update(ctx context.Context, source *models.OfficialSource) (*models.OfficialSource, error) {
	if source.ID == uuid.Nil {
		return nil, exceptions.ErrOfficialSourceNotFound
	}

	result := r.db.WithContext(ctx).
		Model(source).
		Select("title", "url", "agency", "last_verified_at").
		Updates(source)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, exceptions.ErrOfficialSourceNotFound
	}
	return source, nil
}

func (r *officialSourceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return exceptions.ErrOfficialSourceNotFound
	}

	result := r.db.WithContext(ctx).Delete(&models.OfficialSource{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return exceptions.ErrOfficialSourceNotFound
	}
	return nil
}
