package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/carlosEA28/benefitGuide/internal/catalog/internal/models"
	"github.com/carlosEA28/benefitGuide/internal/shared/exceptions"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type requirementService interface {
	CreateRequirement(context.Context, *models.ServiceRequirement) (*models.ServiceRequirement, error)
	ListRequirements(context.Context, uuid.UUID) ([]models.ServiceRequirement, error)
	GetRequirement(context.Context, uuid.UUID, uuid.UUID) (*models.ServiceRequirement, error)
	UpdateRequirementForService(context.Context, uuid.UUID, *models.ServiceRequirement) (*models.ServiceRequirement, error)
	DeleteRequirementForService(context.Context, uuid.UUID, uuid.UUID) error
}

type officialSourceService interface {
	CreateOfficialSource(context.Context, *models.OfficialSource) (*models.OfficialSource, error)
	ListOfficialSources(context.Context, uuid.UUID) ([]models.OfficialSource, error)
	GetOfficialSource(context.Context, uuid.UUID, uuid.UUID) (*models.OfficialSource, error)
	UpdateOfficialSource(context.Context, uuid.UUID, *models.OfficialSource) (*models.OfficialSource, error)
	DeleteOfficialSource(context.Context, uuid.UUID, uuid.UUID) error
}

type catalogService interface {
	requirementService
	officialSourceService
}

type Handler struct {
	requirements *requirementHandler
	sources      *sourceHandler
}

type requirementHandler struct {
	service requirementService
}

func newRequirementHandler(service requirementService) *requirementHandler {
	return &requirementHandler{service: service}
}

func New(service catalogService) *Handler {
	return &Handler{
		requirements: newRequirementHandler(service),
		sources:      newSourceHandler(service),
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	path := "/api/v1/services/{serviceId}/requirements"
	r.Post(path, h.requirements.createRequirement)
	r.Get(path, h.requirements.listRequirements)
	r.Get(path+"/{id}", h.requirements.getRequirement)
	r.Put(path+"/{id}", h.requirements.updateRequirement)
	r.Delete(path+"/{id}", h.requirements.deleteRequirement)

	sourcePath := "/api/v1/services/{serviceId}/sources"
	r.Post(sourcePath, h.sources.createSource)
	r.Get(sourcePath, h.sources.listSources)
	r.Get(sourcePath+"/{id}", h.sources.getSource)
	r.Put(sourcePath+"/{id}", h.sources.updateSource)
	r.Delete(sourcePath+"/{id}", h.sources.deleteSource)
}

func (h *requirementHandler) createRequirement(w http.ResponseWriter, r *http.Request) {
	serviceID, err := parseUUID(r, "serviceId")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	requirement, err := decodeRequirement(r, false)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	requirement.PublicServiceID = serviceID
	created, err := h.service.CreateRequirement(r.Context(), requirement)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, responseFor(created))
}

func (h *requirementHandler) listRequirements(w http.ResponseWriter, r *http.Request) {
	serviceID, err := parseUUID(r, "serviceId")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	requirements, err := h.service.ListRequirements(r.Context(), serviceID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	items := make([]requirementResponse, 0, len(requirements))
	for i := range requirements {
		items = append(items, responseFor(&requirements[i]))
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *requirementHandler) getRequirement(w http.ResponseWriter, r *http.Request) {
	serviceID, err := parseUUID(r, "serviceId")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	id, err := parseUUID(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	requirement, err := h.service.GetRequirement(r.Context(), serviceID, id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, responseFor(requirement))
}

func (h *requirementHandler) updateRequirement(w http.ResponseWriter, r *http.Request) {
	serviceID, err := parseUUID(r, "serviceId")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	id, err := parseUUID(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	changes, err := decodeRequirement(r, true)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	changes.ID = id
	updated, err := h.service.UpdateRequirementForService(r.Context(), serviceID, changes)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, responseFor(updated))
}

func (h *requirementHandler) deleteRequirement(w http.ResponseWriter, r *http.Request) {
	serviceID, err := parseUUID(r, "serviceId")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	id, err := parseUUID(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := h.service.DeleteRequirementForService(r.Context(), serviceID, id); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, exceptions.ErrInvalidRequirement):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid service requirement"})
	case errors.Is(err, exceptions.ErrInvalidOfficialSource):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid official source"})
	case errors.Is(err, exceptions.ErrPublicServiceNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "public service not found"})
	case errors.Is(err, exceptions.ErrRequirementNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "service requirement not found"})
	case errors.Is(err, exceptions.ErrOfficialSourceNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "official source not found"})
	case errors.Is(err, exceptions.ErrPublicServiceInactive):
		writeJSON(w, http.StatusConflict, map[string]string{"error": "public service inactive"})
	default:
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
}
