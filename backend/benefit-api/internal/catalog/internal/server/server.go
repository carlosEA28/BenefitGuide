package server

import (
	"net/http"

	"github.com/carlosEA28/benefitGuide/internal/catalog/internal/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct{ service *service.Service }

func New(service *service.Service) *Handler { return &Handler{service: service} }

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/health", h.healthCheck)
}

func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := h.service.Ready(r.Context()); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"status":"unavailable"}`))
		return
	}
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
