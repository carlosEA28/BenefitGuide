package catalog

import (
	"github.com/carlosEA28/benefitGuide/internal/catalog/internal/repository"
	"github.com/carlosEA28/benefitGuide/internal/catalog/internal/server"
	"github.com/carlosEA28/benefitGuide/internal/catalog/internal/service"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Service = service.Service

type Module struct {
	Repository repository.RequirementRepository
	Service    *service.Service
	Handler    *server.Handler
}

func New(db *gorm.DB) *Module {
	repo := repository.New(db)
	svc := service.New(repo, repo)
	return &Module{Repository: repo, Service: svc, Handler: server.New(svc)}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	m.Handler.RegisterRoutes(r)
}
