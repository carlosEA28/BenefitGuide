package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/carlosEA28/benefitGuide/internal/catalog/internal/models"
	"github.com/google/uuid"
)

type sourceRequest struct {
	Title          *string    `json:"title"`
	URL            *string    `json:"url"`
	Agency         *string    `json:"agency"`
	LastVerifiedAt *time.Time `json:"last_verified_at"`
}

type sourceResponse struct {
	ID              uuid.UUID  `json:"id"`
	PublicServiceID uuid.UUID  `json:"public_service_id"`
	Title           string     `json:"title"`
	URL             string     `json:"url"`
	Agency          string     `json:"agency"`
	LastVerifiedAt  *time.Time `json:"last_verified_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

func decodeSource(r *http.Request) (*models.OfficialSource, error) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var request sourceRequest
	if err := decoder.Decode(&request); err != nil {
		return nil, fmt.Errorf("invalid JSON body: %w", err)
	}
	if request.Title == nil || request.URL == nil || request.Agency == nil {
		return nil, errors.New("missing required source fields")
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return nil, errors.New("JSON body must contain one object")
	}
	return &models.OfficialSource{
		Title:          *request.Title,
		URL:            *request.URL,
		Agency:         *request.Agency,
		LastVerifiedAt: request.LastVerifiedAt,
	}, nil
}

func sourceResponseFor(source *models.OfficialSource) sourceResponse {
	return sourceResponse{
		ID:              source.ID,
		PublicServiceID: source.PublicServiceID,
		Title:           source.Title,
		URL:             source.URL,
		Agency:          source.Agency,
		LastVerifiedAt:  source.LastVerifiedAt,
		CreatedAt:       source.CreatedAt,
	}
}
