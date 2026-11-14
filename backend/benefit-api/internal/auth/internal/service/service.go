package service

import (
	"context"
	"fmt"
	"github.com/carlosEA28/benefitGuide/internal/auth/internal/repository"
)

type Service struct {
	repository *repository.Repository
}

func New(repository *repository.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Ready(ctx context.Context) error {
	if err := s.repository.Ready(ctx); err != nil {
		return fmt.Errorf("database is not ready: %w", err)
	}

	return nil
}
