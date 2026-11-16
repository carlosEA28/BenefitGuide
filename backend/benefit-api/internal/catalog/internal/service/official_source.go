package service

import (
	"context"
	"fmt"
	"net/netip"
	"net/url"
	"strconv"
	"strings"
	"unicode"

	"github.com/carlosEA28/benefitGuide/internal/catalog/internal/models"
	"github.com/carlosEA28/benefitGuide/internal/shared/exceptions"
	"github.com/google/uuid"
)

// CreateOfficialSource validates a source and associates it with an existing public service.
func (s *Service) CreateOfficialSource(ctx context.Context, source *models.OfficialSource) (*models.OfficialSource, error) {
	normalized, err := normalizeOfficialSource(source)
	if err != nil {
		return nil, err
	}
	if _, err := s.publicServiceRepository.FindPublicServiceByID(ctx, normalized.PublicServiceID); err != nil {
		return nil, err
	}
	return s.officialSourceRepository.Create(ctx, normalized)
}

// ListOfficialSources returns the sources of an existing public service.
func (s *Service) ListOfficialSources(ctx context.Context, publicServiceID uuid.UUID) ([]models.OfficialSource, error) {
	if _, err := s.publicServiceRepository.FindPublicServiceByID(ctx, publicServiceID); err != nil {
		return nil, err
	}
	return s.officialSourceRepository.FindByPublicServiceID(ctx, publicServiceID)
}

// GetOfficialSource returns a source only when it belongs to the given public service.
func (s *Service) GetOfficialSource(ctx context.Context, publicServiceID, id uuid.UUID) (*models.OfficialSource, error) {
	if id == uuid.Nil {
		return nil, exceptions.ErrOfficialSourceNotFound
	}
	source, err := s.officialSourceRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if source.PublicServiceID != publicServiceID {
		return nil, exceptions.ErrOfficialSourceNotFound
	}
	return source, nil
}

// UpdateOfficialSource replaces editable fields of a source owned by the given service.
func (s *Service) UpdateOfficialSource(ctx context.Context, publicServiceID uuid.UUID, changes *models.OfficialSource) (*models.OfficialSource, error) {
	if changes == nil {
		return nil, fmt.Errorf("%w: source is nil", exceptions.ErrInvalidOfficialSource)
	}
	current, err := s.GetOfficialSource(ctx, publicServiceID, changes.ID)
	if err != nil {
		return nil, err
	}

	updated := *current
	updated.Title = changes.Title
	updated.URL = changes.URL
	updated.Agency = changes.Agency
	updated.LastVerifiedAt = changes.LastVerifiedAt
	normalized, err := normalizeOfficialSource(&updated)
	if err != nil {
		return nil, err
	}
	return s.officialSourceRepository.Update(ctx, normalized)
}

// DeleteOfficialSource removes a source only under its owning public service.
func (s *Service) DeleteOfficialSource(ctx context.Context, publicServiceID, id uuid.UUID) error {
	if _, err := s.GetOfficialSource(ctx, publicServiceID, id); err != nil {
		return err
	}
	return s.officialSourceRepository.Delete(ctx, id)
}

func normalizeOfficialSource(source *models.OfficialSource) (*models.OfficialSource, error) {
	if source == nil {
		return nil, fmt.Errorf("%w: source is nil", exceptions.ErrInvalidOfficialSource)
	}

	normalized := *source
	normalized.Title = strings.TrimSpace(source.Title)
	normalized.URL = strings.TrimSpace(source.URL)
	if normalized.Title == "" {
		return nil, fmt.Errorf("%w: title is empty", exceptions.ErrInvalidOfficialSource)
	}
	if normalized.URL == "" {
		return nil, fmt.Errorf("%w: URL is empty", exceptions.ErrInvalidOfficialSource)
	}
	if strings.TrimSpace(normalized.Agency) == "" {
		return nil, fmt.Errorf("%w: agency is empty", exceptions.ErrInvalidOfficialSource)
	}
	if !validOfficialSourceURL(normalized.URL) {
		return nil, fmt.Errorf("%w: URL is invalid", exceptions.ErrInvalidOfficialSource)
	}
	return &normalized, nil
}

func validOfficialSourceURL(value string) bool {
	if strings.IndexFunc(value, unicode.IsSpace) >= 0 {
		return false
	}
	parsed, err := url.Parse(value)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return false
	}
	if !validOfficialSourceHost(parsed.Hostname()) {
		return false
	}
	if strings.HasSuffix(parsed.Host, ":") {
		return false
	}
	if port := parsed.Port(); port != "" {
		portNumber, err := strconv.Atoi(port)
		if err != nil || portNumber < 1 || portNumber > 65535 {
			return false
		}
	}
	return true
}

func validOfficialSourceHost(host string) bool {
	if host == "" {
		return false
	}
	if _, err := netip.ParseAddr(host); err == nil {
		return true
	}
	if len(host) > 253 {
		return false
	}
	for _, label := range strings.Split(host, ".") {
		if len(label) == 0 || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}
		for _, character := range label {
			if character != '-' && (character < 'a' || character > 'z') && (character < 'A' || character > 'Z') && (character < '0' || character > '9') {
				return false
			}
		}
	}
	return true
}
