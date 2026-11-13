package service

import (
	"context"
	"fmt"
	authmodule "github.com/carlosEA28/benefitGuide/internal/auth"
	catalogmodule "github.com/carlosEA28/benefitGuide/internal/catalog"
	"github.com/carlosEA28/benefitGuide/internal/consultations/internal/repository"
)

type Service struct {
	repository *repository.Repository
	auth       *authmodule.Service
	catalog    *catalogmodule.Service
}

func New(repository *repository.Repository, auth *authmodule.Service, catalog *catalogmodule.Service) *Service {
	return &Service{repository: repository, auth: auth, catalog: catalog}
}

func (s *Service) Ready(ctx context.Context) error {
	if err := s.repository.Ready(ctx); err != nil {
		return fmt.Errorf("database is not ready: %w", err)
	}
	if err := s.auth.Ready(ctx); err != nil {
		return fmt.Errorf("auth is not ready: %w", err)
	}
	if err := s.catalog.Ready(ctx); err != nil {
		return fmt.Errorf("catalog is not ready: %w", err)
	}
	return nil
}
