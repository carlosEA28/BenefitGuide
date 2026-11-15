package service

import (
	"context"
	"fmt"

	"github.com/carlosEA28/benefitGuide/internal/consultations/internal/repository"
)

// Service holds the service dependencies and provides business logic.
// The repository must implement the repository.Requirement interface to define
// access requirements specific to this service.
type Service struct {
	repository repository.Requirement
}

// New creates a new Service instance with a requirement-aware repository.
func New(repository repository.Requirement) *Service {
	return &Service{
		repository: repository,
	}
}

// Ready checks if the database is healthy and all requirements are met.
// It first checks the database connection, then validates any service-specific
// requirements defined in the repository's Requirement interface.
func (s *Service) Ready(ctx context.Context) error {
	if err := s.repository.Check(ctx); err != nil {
		return fmt.Errorf("requirement check failed: %w", err)
	}

	return nil
}
