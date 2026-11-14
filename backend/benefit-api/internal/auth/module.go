package auth

import (
	"github.com/carlosEA28/benefitGuide/internal/auth/internal/repository"
	"github.com/carlosEA28/benefitGuide/internal/auth/internal/server"
	"github.com/carlosEA28/benefitGuide/internal/auth/internal/service"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Service = service.Service

type Module struct {
	Repository *repository.Repository
	Service    *service.Service
	Handler    *server.Handler
}

func New(db *gorm.DB) *Module {
	repo := repository.New(db)
	svc := service.New(repo)
	return &Module{Repository: repo, Service: svc, Handler: server.New(svc)}
}

func (m *Module) RegisterRoutes(r chi.Router) {
	r.Route("/auth", m.Handler.RegisterRoutes)
}
