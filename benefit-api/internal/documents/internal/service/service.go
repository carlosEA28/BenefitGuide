package service

import (
	"context"
	"fmt"
	consultationsmodule "github.com/carlosEA28/benefitGuide/internal/consultations"
	"github.com/carlosEA28/benefitGuide/internal/documents/internal/repository"
)

type Service struct {
	repository    *repository.Repository
	consultations *consultationsmodule.Service
}

func New(repository *repository.Repository, consultations *consultationsmodule.Service) *Service {
	return &Service{repository: repository, consultations: consultations}
}

func (s *Service) Ready(ctx context.Context) error {
	if err := s.repository.Ready(ctx); err != nil {
		return fmt.Errorf("database is not ready: %w", err)
	}
	if err := s.consultations.Ready(ctx); err != nil {
		return fmt.Errorf("consultations is not ready: %w", err)
	}
	return nil
}
