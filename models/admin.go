package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Admin model struct
type Admin struct {
	gorm.Model `json:"-"`
	ID         uint `gorm:"column:id;not null;unique; primaryKey;" json:"expert_id"`
	//Admin permissions
	CanManageExpert  bool `gorm:"column:can_manage_expert" json:"can_manage_expert"`
	CanManageLearner bool `gorm:"column:can_manage_learner" json:"can_manage_learner"`
	CanManageAdmin   bool `gorm:"column:can_manage_admin" json:"can_manage_admin"`
	//An admin can have only one account
	AccountID uuid.UUID `gorm:"size:225;column:account_id"`
	Account   Account   `gorm:"foreignKey:AccountID"`
	//default timestamps
	CreatedAt time.Time  `gorm:"column:CreatedAt;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:UpdatedAt;autoCreateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:DeletedAt" sql:"index" json:"deleted_at";`
}
