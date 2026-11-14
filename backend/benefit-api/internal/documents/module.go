package documents

import (
	consultationsmodule "github.com/carlosEA28/benefitGuide/internal/consultations"
	"github.com/carlosEA28/benefitGuide/internal/documents/internal/repository"
	"github.com/carlosEA28/benefitGuide/internal/documents/internal/server"
	"github.com/carlosEA28/benefitGuide/internal/documents/internal/service"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Service = service.Service

type Module struct {
	Repository *repository.Repository
	Service    *service.Service
	Handler    *server.Handler
}

func New(db *gorm.DB, consultations *consultationsmodule.Service) *Module {
	repo := repository.New(db)
	svc := service.New(repo, consultations)
	return &Module{Repository: repo, Service: svc, Handler: server.New(svc)}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/documents", m.Handler.RegisterRoutes)
}
