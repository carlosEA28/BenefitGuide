package models

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	ServiceStatusActive   Status = "ACTIVE"
	ServiceStatusInactive Status = "INACTIVE"
)

type PublicService struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	ExternalID   *string `gorm:"column:external_id;type:varchar(150)" json:"external_id,omitempty"`
	SourceSystem string  `gorm:"column:source_system;type:varchar(50);not null" json:"source_system" validate:"required,max=50"`

	Name        string `gorm:"type:varchar(255);not null" json:"name" validate:"required,max=255"`
	Description string `gorm:"type:text" json:"description"`
	Category    string `gorm:"type:varchar(100)" json:"category"`
	Agency      string `gorm:"type:varchar(255)" json:"agency"`

	Status Status `gorm:"type:varchar(30);not null;default:ACTIVE" json:"status" validate:"required,oneof=ACTIVE INACTIVE"`

	CreatedAt    time.Time  `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"not null;autoUpdateTime" json:"updated_at"`
	LastSyncedAt *time.Time `gorm:"column:last_synced_at" json:"last_synced_at,omitempty"`
}
