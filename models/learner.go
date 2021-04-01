package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Learner model struct
type Learner struct {
	gorm.Model `json:"-"`
	ID         uint `gorm:"column:id;not null;unique; primaryKey;" json:"ID"`
	//A learner can have only one account
	AccountID          uuid.UUID `gorm:"size:225;column:account_id"`
	Account            Account   `gorm:"foreignKey:AccountID" json:"-"`
	AvailableCoinCount uint
	TranslationSession []*TranslationSession `gorm:"many2many:translation_session_learners;"`

	//default timestamps
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
