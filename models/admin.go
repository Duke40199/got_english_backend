package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Admin model struct
type Admin struct {
	gorm.Model
	ID uint `gorm:"column:id;not null;unique; primaryKey;" json:"expert_id"`

	//Admin permissions
	CanManageExpert  bool `gorm:"column:can_manage_expert" json:"can_manage_expert"`
	CanManageLearner bool `gorm:"column:can_manage_learner" json:"can_manage_learner"`
	CanManageAdmin   bool `gorm:"column:can_manage_admin" json:"can_manage_admin"`
	//An admin can have only one account
	AccountsID uuid.UUID `gorm:"size:225;column:accounts_id"`
	//default timestamps
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime"`
}
