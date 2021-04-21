package models

import (
	"time"

	"github.com/google/uuid"
)

// Learner model struct
type Learner struct {
	ID                 uint                  `gorm:"column:id;not null;unique; primaryKey;" json:"id"`
	AvailableCoinCount uint                  `gorm:"column:available_coin_count" json:"available_coin_count"`
	TranslationSession []*TranslationSession `gorm:"many2many:translation_session_learners;" json:"-"`
	//A learner can have only one account
	AccountID uuid.UUID `gorm:"column:account_id" json:"account_id"`
	Account   Account   `gorm:"foreignKey:AccountID" json:"account"`

	//default timestamps
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
}
