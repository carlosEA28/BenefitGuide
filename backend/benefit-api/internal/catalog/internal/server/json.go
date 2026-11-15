package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/carlosEA28/benefitGuide/internal/catalog/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type requirementRequest struct {
	Description          *string                 `json:"description"`
	Type                 *models.RequirementType `json:"type"`
	ConditionDescription *string                 `json:"condition_description"`
	DisplayOrder         *int                    `json:"display_order"`
}

type requirementResponse struct {
	ID                   uuid.UUID              `json:"id"`
	PublicServiceID      uuid.UUID              `json:"public_service_id"`
	Description          string                 `json:"description"`
	Type                 models.RequirementType `json:"type"`
	ConditionDescription *string                `json:"condition_description"`
	DisplayOrder         int                    `json:"display_order"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
}

func parseUUID(r *http.Request, name string) (uuid.UUID, error) {
	id, err := uuid.Parse(chi.URLParam(r, name))
	if err != nil || id == uuid.Nil {
		return uuid.Nil, fmt.Errorf("invalid %s", name)
	}
	return id, nil
}

func decodeRequirement(r *http.Request, requireOrder bool) (*models.ServiceRequirement, error) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var request requirementRequest
	if err := decoder.Decode(&request); err != nil {
		return nil, fmt.Errorf("invalid JSON body: %w", err)
	}
	if request.Description == nil || request.Type == nil || (requireOrder && request.DisplayOrder == nil) {
		return nil, errors.New("missing required requirement fields")
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return nil, errors.New("JSON body must contain one object")
	}
	model := &models.ServiceRequirement{
		Description:          *request.Description,
		Type:                 *request.Type,
		ConditionDescription: request.ConditionDescription,
	}
	if request.DisplayOrder != nil {
		model.DisplayOrder = *request.DisplayOrder
	}
	return model, nil
}

func responseFor(requirement *models.ServiceRequirement) requirementResponse {
	return requirementResponse{
		ID:                   requirement.ID,
		PublicServiceID:      requirement.PublicServiceID,
		Description:          requirement.Description,
		Type:                 requirement.Type,
		ConditionDescription: requirement.ConditionDescription,
		DisplayOrder:         requirement.DisplayOrder,
		CreatedAt:            requirement.CreatedAt,
		UpdatedAt:            requirement.UpdatedAt,
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	encoded, err := json.Marshal(body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("{\"error\":\"internal server error\"}\n"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(append(encoded, '\n'))
}
