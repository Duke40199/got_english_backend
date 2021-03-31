package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PrivateCallSession model struct
type PrivateCallSession struct {
	gorm.Model `json:"-"`
	ID         uint      `gorm:"column:id;autoIncrement;not null; unique; primaryKey;" json:"id"`
	Duration   time.Time `json:"duration"`
	StartedAt  time.Time `gorm:"column:started_at" json:"started_at"`
	FinishedAt time.Time `gorm:"column:finished_at" json:"finished_at"`
	//Rating attributes
	Rating            float32
	RatingDescription string
	//A private call session can have only one learner
	Learners  Learner   `gorm:"foreignKey:LearnerID"`
	LearnerID uuid.UUID `gorm:"size:255"`
	//A private call session can have only one expert
	Expert   Expert    `gorm:"foreignKey:ExpertID"`
	ExpertID uuid.UUID `gorm:"size:255"`
	//An messaging session can only have one pricing.
	Pricing   Pricing `gorm:"foreignKey:PricingID"`
	PricingID uint    `gorm:"size:255"`
	//default timestamps
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
