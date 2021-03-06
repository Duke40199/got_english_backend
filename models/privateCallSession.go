package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PrivateCallSession model struct
type PrivateCallSession struct {
	gorm.Model
	ID         uint      `gorm:"column:id;autoIncrement;not null; unique; primaryKey;" json:"id"`
	Duration   time.Time `json:"duration"`
	StartedAt  time.Time `gorm:"" json:"started_at"`
	FinishedAt time.Time `gorm:"" json:"finished_at"`
	//Rating attributes
	Rating            float32
	RatingDescription string

	//A private call session can have only one learner
	Learners  Learner   `gorm:"foreignKey:LearnerID"`
	LearnerID uuid.UUID `gorm:"size:255"`
	//A private call session can have only one expert
	Expert   Expert    `gorm:"foreignKey:ExpertID"`
	ExpertID uuid.UUID `gorm:"size:255"`
	//A PrivateCall session will have a pricing
	Price float32 `gorm:"column:price"`
}
