package models

import (
	"time"

	"github.com/google/uuid"
)

type OfficialSource struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	PublicServiceID uuid.UUID `gorm:"column:public_service_id;type:uuid;not null;index" json:"public_service_id" validate:"required"`

	Title  string `gorm:"type:varchar(255);not null" json:"title" validate:"required,max=255"`
	URL    string `gorm:"column:url;type:text;not null" json:"url" validate:"required,url"`
	Agency string `gorm:"type:varchar(255);not null" json:"agency"`

	LastVerifiedAt *time.Time `gorm:"column:last_verified_at" json:"last_verified_at,omitempty"`
	CreatedAt      time.Time  `gorm:"not null;autoCreateTime" json:"created_at"`
}
