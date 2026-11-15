package service

import (
	"github.com/carlosEA28/benefitGuide/internal/catalog/internal/repository"
)

// Service holds the service dependencies and provides catalog business logic.
type Service struct {
	repository              repository.RequirementRepository
	publicServiceRepository repository.PublicServiceRepository
}

// New creates a catalog Service with requirement and public-service repositories.
func New(requirements repository.RequirementRepository, publicServices repository.PublicServiceRepository) *Service {
	return &Service{
		repository:              requirements,
		publicServiceRepository: publicServices,
	}
}
