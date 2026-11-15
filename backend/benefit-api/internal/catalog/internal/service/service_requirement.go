package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/carlosEA28/benefitGuide/internal/catalog/internal/models"
	"github.com/carlosEA28/benefitGuide/internal/shared/exceptions"
	"github.com/google/uuid"
)

// CreateRequirement validates and creates a requirement for an active public service.
func (s *Service) CreateRequirement(ctx context.Context, requirement *models.ServiceRequirement) (*models.ServiceRequirement, error) {
	if err := validateRequirement(requirement); err != nil {
		return nil, err
	}

	publicService, err := s.publicServiceRepository.FindPublicServiceByID(ctx, requirement.PublicServiceID)
	if err != nil {
		return nil, err
	}
	if publicService.Status != models.ServiceStatusActive {
		return nil, exceptions.ErrPublicServiceInactive
	}

	return s.repository.CreateRequirement(ctx, requirement)
}

// UpdateRequirement updates mutable fields of an existing requirement.
func (s *Service) UpdateRequirement(ctx context.Context, changes *models.ServiceRequirement) (*models.ServiceRequirement, error) {
	if changes == nil {
		return nil, fmt.Errorf("%w: requirement is nil", exceptions.ErrInvalidRequirement)
	}
	if changes.ID == uuid.Nil {
		return nil, exceptions.ErrRequirementNotFound
	}

	current, err := s.repository.FindRequirementByID(ctx, changes.ID)
	if err != nil {
		return nil, err
	}

	updated := *current
	updated.Description = changes.Description
	updated.Type = changes.Type
	updated.ConditionDescription = changes.ConditionDescription
	updated.DisplayOrder = changes.DisplayOrder
	if updated.Type == models.RequirementMandatory {
		updated.ConditionDescription = nil
	}
	if err := validateRequirement(&updated); err != nil {
		return nil, err
	}
	return s.repository.UpdateRequirement(ctx, &updated)
}

// DeleteRequirement removes an existing requirement.
func (s *Service) DeleteRequirement(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return exceptions.ErrRequirementNotFound
	}
	return s.repository.DeleteRequirement(ctx, id)
}

func validateRequirement(requirement *models.ServiceRequirement) error {
	if requirement == nil {
		return fmt.Errorf("%w: requirement is nil", exceptions.ErrInvalidRequirement)
	}
	if strings.TrimSpace(requirement.Description) == "" {
		return fmt.Errorf("%w: description is empty", exceptions.ErrInvalidRequirement)
	}
	if requirement.Type != models.RequirementMandatory && requirement.Type != models.RequirementConditional {
		return fmt.Errorf("%w: unsupported type", exceptions.ErrInvalidRequirement)
	}
	if requirement.Type == models.RequirementConditional &&
		(requirement.ConditionDescription == nil || strings.TrimSpace(*requirement.ConditionDescription) == "") {
		return fmt.Errorf("%w: condition description is empty", exceptions.ErrInvalidRequirement)
	}
	return nil
}
