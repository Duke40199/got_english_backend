package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TranslationSession model struct
type TranslationSession struct {
	gorm.Model        `json:"-"`
	ID                uint      `gorm:"column:id;autoIncrement;not null; unique; primaryKey;" json:"id"`
	DurationInSeconds uint      `gorm:"column:duration_in_seconds" json:"duration"`
	StartedAt         time.Time `gorm:"" json:"started_at"`
	FinishedAt        time.Time `gorm:"" json:"finished_at"`
	//A translation session can have many learners
	Learners []*Learner `gorm:"many2many:translation_session_learners;"`
	//A translation session can have only one expert.
	Expert   Expert    `gorm:"foreignKey:ExpertID"`
	ExpertID uuid.UUID `gorm:"size:255"`
	//An messaging session can only have one pricing.
	Pricing   Pricing `gorm:"foreignKey:PricingID"`
	PricingID uint    `gorm:"size:255"`
	//A tranlsation session can have only one rating
	Rating   Rating `gorm:"foreignKey:RatingID"`
	RatingID uint   `gorm:"size:255"`
	//default timestamps
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
}
