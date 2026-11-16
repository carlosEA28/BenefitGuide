package service

import (
	"github.com/carlosEA28/benefitGuide/internal/catalog/internal/repository"
)

// Service holds the service dependencies and provides catalog business logic.
type Service struct {
	repository               repository.RequirementRepository
	publicServiceRepository  repository.PublicServiceRepository
	officialSourceRepository repository.OfficialSourceRepository
}

// New creates a catalog Service with its repositories.
func New(requirements repository.RequirementRepository, publicServices repository.PublicServiceRepository, officialSources ...repository.OfficialSourceRepository) *Service {
	service := &Service{
		repository:              requirements,
		publicServiceRepository: publicServices,
	}
	if len(officialSources) > 0 {
		service.officialSourceRepository = officialSources[0]
	}
	return service
}
