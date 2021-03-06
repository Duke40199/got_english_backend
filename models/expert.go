package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Expert model struct
type Expert struct {
	gorm.Model
	ID uint `gorm:"column:id;not null;unique; primaryKey;" json:"expert_id"`
	//Expert permissions
	CanChat                   bool `gorm:"column:can_chat" json:"can_chat"`
	CanJoinTranslationSession bool `gorm:"column:can_join_translation_session" json:"can_join_translation_session"`
	CanJoinPrivateCallSession bool `gorm:"column:can_private_call_session" json:"can_private_call_session"`
	//An expert can have only one account
	AccountsID uuid.UUID `gorm:"size:225;column:accounts_id"`
	//default timestamps
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime"`
}
