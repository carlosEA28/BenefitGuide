package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type healthService interface {
	Ready(context.Context) error
}

type Handler struct{ service healthService }

func New(service healthService) *Handler { return &Handler{service: service} }

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
