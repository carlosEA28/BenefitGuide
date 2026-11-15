package catalog

import (
	"github.com/carlosEA28/benefitGuide/internal/catalog/internal/repository"
	"github.com/carlosEA28/benefitGuide/internal/catalog/internal/service"
	"gorm.io/gorm"
)

type Service = service.Service

type Module struct {
	Repository repository.RequirementRepository
	Service    *service.Service
}

func New(db *gorm.DB) *Module {
	repo := repository.New(db)
	svc := service.New(repo, repo)
	return &Module{Repository: repo, Service: svc}
}
