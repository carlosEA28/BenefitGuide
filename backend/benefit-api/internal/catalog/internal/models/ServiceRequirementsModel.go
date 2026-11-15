package models

import (
	"time"

	"github.com/google/uuid"
)

type RequirementType string

const (
	RequirementMandatory   RequirementType = "MANDATORY"
	RequirementConditional RequirementType = "CONDITIONAL"
)

type ServiceRequirement struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`

	PublicServiceID uuid.UUID `gorm:"column:public_service_id;type:uuid;not null;index" json:"public_service_id" validate:"required"`

	Description string `gorm:"type:text;not null" json:"description" validate:"required"`

	Type RequirementType `gorm:"type:varchar(30);not null" json:"type" validate:"required,oneof=MANDATORY CONDITIONAL"`

	ConditionDescription *string `gorm:"type:text" json:"condition_description,omitempty"`

	DisplayOrder int `gorm:"column:display_order;not null;default:0" json:"display_order"`

	CreatedAt time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`
}
