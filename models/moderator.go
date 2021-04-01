package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Moderator model struct
type Moderator struct {
	gorm.Model `json:"-"`
	ID         uint `gorm:"column:id;not null;unique; primaryKey;" json:"id"`

	//Moderator permissions
	CanManageCoinBundle      bool `gorm:"column:can_manage_coin_bundle" json:"can_manage_coin_bundle"`
	CanManagePricing         bool `gorm:"column:can_manage_pricing" json:"can_manage_pricing"`
	CanManageApplicationForm bool `gorm:"column:can_manage_application_form" json:"can_manage_application_form"`
	//An expert can have only one account
	AccountID uuid.UUID `gorm:"column:account_id" json:"account_id"`

	//default timestamps
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
