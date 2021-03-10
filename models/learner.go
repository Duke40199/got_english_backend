package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Learner model struct
type Learner struct {
	gorm.Model
	ID uint `gorm:"column:id;not null;unique; primaryKey;" json:"ID"`
	//A learner can have only one account
	AccountsID         uuid.UUID             `gorm:"size:225;column:accounts_id"`
	TranslationSession []*TranslationSession `gorm:"many2many:translation_session_learners;"`
	//default timestamps
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime"`
}
