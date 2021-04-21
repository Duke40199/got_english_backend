package models

import (
	"time"

	"github.com/google/uuid"
)

// Admin model struct
type Admin struct {
	ID uint `gorm:"column:id;not null;unique; primaryKey;" json:"id"`
	//Admin permissions
	CanManageExpert    bool `gorm:"column:can_manage_expert" json:"can_manage_expert"`
	CanManageLearner   bool `gorm:"column:can_manage_learner" json:"can_manage_learner"`
	CanManageAdmin     bool `gorm:"column:can_manage_admin" json:"can_manage_admin"`
	CanManageModerator bool `gorm:"column:can_manage_moderator" json:"can_manage_moderator"`
	//An admin can have only one account
	AccountID uuid.UUID `gorm:"column:account_id" json:"account_id"`
	//default timestamps
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
}
