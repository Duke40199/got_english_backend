package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Expert model struct
type Expert struct {
	gorm.Model `json:"-"`
	ID         uint   `gorm:"column:id;not null;unique; primaryKey;" json:"expert_id"`
	Profession string `gorm:"column:profession" json:"professtion"`
	//Expert permissions
	CanChat                   bool `gorm:"column:can_chat" json:"can_chat;omitempty"`
	CanJoinTranslationSession bool `gorm:"column:can_join_translation_session" json:"can_join_translation_session;omitempty"`
	CanJoinPrivateCallSession bool `gorm:"column:can_join_private_call_session" json:"can_join_private_call_session;omitempty"`
	//An expert can have only one account
	AccountID uuid.UUID `gorm:"size:225;column:account_id;unique"`
	Account   Account   `gorm:"foreignKey:AccountID;"`
	//default timestamps
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
