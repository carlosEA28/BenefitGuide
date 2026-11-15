package service

import (
	"context"
	"errors"
	"testing"

	"github.com/carlosEA28/benefitGuide/internal/catalog/internal/models"
	"github.com/carlosEA28/benefitGuide/internal/shared/exceptions"
	"github.com/google/uuid"
)

type requirementRepoStub struct {
	created int
	updated int
	deleted int
	records map[uuid.UUID]models.ServiceRequirement
}

func (r *requirementRepoStub) CreateRequirement(_ context.Context, requirement *models.ServiceRequirement) (*models.ServiceRequirement, error) {
	r.created++
	return requirement, nil
}
func (r *requirementRepoStub) FindRequirementByID(_ context.Context, id uuid.UUID) (*models.ServiceRequirement, error) {
	requirement, ok := r.records[id]
	if !ok {
		return nil, exceptions.ErrRequirementNotFound
	}
	return &requirement, nil
}
func (r *requirementRepoStub) FindByPublicServiceID(context.Context, uuid.UUID) ([]models.ServiceRequirement, error) {
	panic("not used")
}
func (r *requirementRepoStub) UpdateRequirement(_ context.Context, requirement *models.ServiceRequirement) (*models.ServiceRequirement, error) {
	r.updated++
	if _, ok := r.records[requirement.ID]; !ok {
		return nil, exceptions.ErrRequirementNotFound
	}
	r.records[requirement.ID] = *requirement
	return requirement, nil
}
func (r *requirementRepoStub) DeleteRequirement(_ context.Context, id uuid.UUID) error {
	r.deleted++
	if _, ok := r.records[id]; !ok {
		return exceptions.ErrRequirementNotFound
	}
	delete(r.records, id)
	return nil
}

type publicServiceRepoStub struct {
	service *models.PublicService
	err     error
	lookups int
}

func (r *publicServiceRepoStub) FindPublicServiceByID(_ context.Context, _ uuid.UUID) (*models.PublicService, error) {
	r.lookups++
	return r.service, r.err
}

func TestCreateRequirement(t *testing.T) {
	condition := "apenas se empregado"
	serviceID := uuid.New()
	valid := func() *models.ServiceRequirement {
		return &models.ServiceRequirement{
			PublicServiceID:      serviceID,
			Description:          "comprovar renda",
			Type:                 models.RequirementConditional,
			ConditionDescription: &condition,
		}
	}

	tests := []struct {
		name          string
		requirement   *models.ServiceRequirement
		publicService *models.PublicService
		publicErr     error
		wantErr       error
		wantCreate    int
		wantLookup    int
	}{
		{name: "active service", requirement: valid(), publicService: &models.PublicService{ID: serviceID, Status: models.ServiceStatusActive}, wantCreate: 1, wantLookup: 1},
		{name: "inactive service", requirement: valid(), publicService: &models.PublicService{ID: serviceID, Status: models.ServiceStatusInactive}, wantErr: exceptions.ErrPublicServiceInactive, wantLookup: 1},
		{name: "missing service", requirement: valid(), publicErr: exceptions.ErrPublicServiceNotFound, wantErr: exceptions.ErrPublicServiceNotFound, wantLookup: 1},
		{name: "blank description", requirement: func() *models.ServiceRequirement { r := valid(); r.Description = " \t "; return r }(), wantErr: exceptions.ErrInvalidRequirement},
		{name: "invalid type", requirement: func() *models.ServiceRequirement { r := valid(); r.Type = "OTHER"; return r }(), wantErr: exceptions.ErrInvalidRequirement},
		{name: "empty condition", requirement: func() *models.ServiceRequirement {
			r := valid()
			empty := " \t "
			r.ConditionDescription = &empty
			return r
		}(), wantErr: exceptions.ErrInvalidRequirement},
		{name: "missing condition", requirement: func() *models.ServiceRequirement { r := valid(); r.ConditionDescription = nil; return r }(), wantErr: exceptions.ErrInvalidRequirement},
		{name: "nil requirement", wantErr: exceptions.ErrInvalidRequirement},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			requirements := &requirementRepoStub{}
			services := &publicServiceRepoStub{service: tc.publicService, err: tc.publicErr}
			got, err := New(requirements, services).CreateRequirement(context.Background(), tc.requirement)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("error = %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr == nil && got != tc.requirement {
				t.Fatalf("created requirement = %#v", got)
			}
			if requirements.created != tc.wantCreate || services.lookups != tc.wantLookup {
				t.Fatalf("creates = %d, lookups = %d; want %d, %d", requirements.created, services.lookups, tc.wantCreate, tc.wantLookup)
			}
		})
	}
}

func TestUpdateRequirementPreservesIdentityAndClearsCondition(t *testing.T) {
	id := uuid.New()
	serviceID := uuid.New()
	condition := "se houver renda"
	requirements := &requirementRepoStub{records: map[uuid.UUID]models.ServiceRequirement{
		id: {
			ID: id, PublicServiceID: serviceID, Description: "renda",
			Type: models.RequirementConditional, ConditionDescription: &condition,
		},
	}}
	services := &publicServiceRepoStub{}
	changes := &models.ServiceRequirement{
		ID: id, PublicServiceID: uuid.New(), Description: "nova descrição",
		Type: models.RequirementMandatory, ConditionDescription: &condition, DisplayOrder: 2,
	}

	got, err := New(requirements, services).UpdateRequirement(context.Background(), changes)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != id || got.PublicServiceID != serviceID || got.ConditionDescription != nil || got.Description != changes.Description || got.DisplayOrder != 2 {
		t.Fatalf("updated requirement = %#v", got)
	}
	if requirements.updated != 1 || services.lookups != 0 || requirements.records[id].PublicServiceID != serviceID {
		t.Fatalf("repository state = %#v, updates = %d, service lookups = %d", requirements.records, requirements.updated, services.lookups)
	}
}

func TestUpdateRequirementRejectsMissingAndInvalid(t *testing.T) {
	id := uuid.New()
	requirements := &requirementRepoStub{records: map[uuid.UUID]models.ServiceRequirement{}}
	svc := New(requirements, &publicServiceRepoStub{})
	if _, err := svc.UpdateRequirement(context.Background(), &models.ServiceRequirement{ID: id}); !errors.Is(err, exceptions.ErrRequirementNotFound) {
		t.Fatalf("missing requirement error = %v", err)
	}
	if _, err := svc.UpdateRequirement(context.Background(), &models.ServiceRequirement{}); !errors.Is(err, exceptions.ErrRequirementNotFound) {
		t.Fatalf("nil ID error = %v", err)
	}
	if requirements.updated != 0 {
		t.Fatalf("unexpected updates: %d", requirements.updated)
	}

	requirements.records[id] = models.ServiceRequirement{ID: id, PublicServiceID: uuid.New(), Description: "valid", Type: models.RequirementMandatory}
	for _, changes := range []*models.ServiceRequirement{
		{ID: id, Description: " ", Type: models.RequirementMandatory},
		{ID: id, Description: "valid", Type: "OTHER"},
		{ID: id, Description: "valid", Type: models.RequirementConditional},
	} {
		if _, err := svc.UpdateRequirement(context.Background(), changes); !errors.Is(err, exceptions.ErrInvalidRequirement) {
			t.Fatalf("invalid update error = %v", err)
		}
	}
	if requirements.updated != 0 {
		t.Fatalf("unexpected updates: %d", requirements.updated)
	}
}

func TestDeleteRequirement(t *testing.T) {
	id := uuid.New()
	requirements := &requirementRepoStub{records: map[uuid.UUID]models.ServiceRequirement{
		id: {ID: id},
	}}
	svc := New(requirements, &publicServiceRepoStub{})
	if err := svc.DeleteRequirement(context.Background(), id); err != nil {
		t.Fatal(err)
	}
	if err := svc.DeleteRequirement(context.Background(), id); !errors.Is(err, exceptions.ErrRequirementNotFound) {
		t.Fatalf("second delete error = %v", err)
	}
	if err := svc.DeleteRequirement(context.Background(), uuid.Nil); !errors.Is(err, exceptions.ErrRequirementNotFound) {
		t.Fatalf("nil ID error = %v", err)
	}
	if requirements.deleted != 2 {
		t.Fatalf("repository deletes = %d, want 2", requirements.deleted)
	}
}
