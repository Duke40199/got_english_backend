package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Moderator model struct
type Moderator struct {
	gorm.Model
	ID uint `gorm:"column:id;not null;unique; primaryKey;" json:"expert_id"`

	//Moderator permissions
	CanManageCoinBundle      bool `gorm:"column:can_manage_coin_bundle" json:"can_manage_coin_bundle"`
	CanManagePricing         bool `gorm:"column:can_manage_pricing" json:"can_manage_pricing"`
	CanManageApplicationForm bool `gorm:"column:can_manage_application_form" json:"can_manage_application_form"`
	//An expert can have only one user account
	AccountsID uuid.UUID `gorm:"size:225;column:accounts_id"`
	//default timestamps
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime"`
}
