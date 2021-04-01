package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Learner model struct
type Learner struct {
	gorm.Model `json:"-"`
	ID         uint `gorm:"column:id;not null;unique; primaryKey;" json:"id"`

	//A learner can have only one account
	AccountID uuid.UUID `gorm:"column:account_id" json:"account_id"`

	AvailableCoinCount uint                  `gorm:"column:available_coin_count" json:"available_coin_count"`
	TranslationSession []*TranslationSession `gorm:"many2many:translation_session_learners;" json:"-"`

	//default timestamps
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
