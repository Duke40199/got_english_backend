package models

import (
	"time"

	"github.com/google/uuid"
)

// Expert model struct
type Expert struct {
	ID            uint    `gorm:"column:id;not null;unique; primaryKey;" json:"id"`
	AverageRating float32 `gorm:"-" json:"average_rating"`
	//Expert permissions
	CanChat                   bool    `gorm:"column:can_chat" json:"can_chat"`
	CanJoinTranslationSession bool    `gorm:"column:can_join_translation_session" json:"can_join_translation_session"`
	CanJoinLiveCallSession    bool    `gorm:"column:can_join_live_call_session" json:"can_join_live_call_session"`
	WeightedRating            float32 `gorm:"column:weighted_rating;default:0" json:"weighted_rating"`
	//An expert can have only one account
	Account            *Account              `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	AccountID          uuid.UUID             `gorm:"column:account_id" json:"account_id"`
	TranslationSession *[]TranslationSession `gorm:"-" json:"translation_session,omitempty"`
	//default timestamps
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
}
