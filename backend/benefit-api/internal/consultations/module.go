package consultations

import (
	authmodule "github.com/carlosEA28/benefitGuide/internal/auth"
	catalogmodule "github.com/carlosEA28/benefitGuide/internal/catalog"
	"github.com/carlosEA28/benefitGuide/internal/consultations/internal/repository"
	"github.com/carlosEA28/benefitGuide/internal/consultations/internal/server"
	"github.com/carlosEA28/benefitGuide/internal/consultations/internal/service"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Service = service.Service

type Module struct {
	Repository *repository.Repository
	Service    *service.Service
	Handler    *server.Handler
}

func New(db *gorm.DB, auth *authmodule.Service, catalog *catalogmodule.Service) *Module {
	repo := repository.New(db)
	svc := service.New(repo, auth, catalog)
	return &Module{Repository: repo, Service: svc, Handler: server.New(svc)}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/consultations", m.Handler.RegisterRoutes)
}
