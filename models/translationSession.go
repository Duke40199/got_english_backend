package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TranslationSession model struct
type TranslationSession struct {
	gorm.Model
	ID         uint      `gorm:"column:id;autoIncrement;not null; unique; primaryKey;" json:"id"`
	Duration   time.Time `json:"duration"`
	StartedAt  time.Time `gorm:"" json:"started_at"`
	FinishedAt time.Time `gorm:"" json:"finished_at"`
	//Rating attributes
	Rating            float32
	RatingDescription string

	//A translation session can have many learners
	Learners []*Learner `gorm:"many2many:translation_session_learners;"`
	//A translation session can have only one expert.
	Expert   Expert    `gorm:"foreignKey:ExpertID"`
	ExpertID uuid.UUID `gorm:"size:255"`
	//default timestamps
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime"`
}
