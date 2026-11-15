package service

import (
	"github.com/carlosEA28/benefitGuide/internal/catalog/internal/repository"
)

// Service holds the service dependencies and provides catalog business logic.
type Service struct {
	repository repository.RequirementRepository
}

// New creates a new Service instance with a requirement-aware repository.
func New(repository repository.RequirementRepository) *Service {
	return &Service{
		repository: repository,
	}
}
