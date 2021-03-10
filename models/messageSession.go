package models

import (
	"time"

	"gorm.io/gorm"
)

// MessageSession model struct
type MessageSession struct {
	gorm.Model
	ID uint `gorm:"column:id;autoIncrement;not null; unique; primaryKey;" json:"id"`
	//A leaner can have many message sessions.
	Learner   Learner `gorm:"foreignKey:LearnerID"`
	LearnerID uint    `gorm:"size:255"`
	//An expert can have many message sessions.
	Expert   Expert `gorm:"foreignKey:ExpertID"`
	ExpertID uint   `gorm:"size:255"`
	//MessageSession status
	IsCancelled bool      `gorm:"default:false;" json:"is_cancelled"`
	IsFinished  bool      `gorm:"default:false;" json:"is_finished"`
	StartedAt   time.Time `gorm:"" json:"started_at"`
	FinishedAt  time.Time `gorm:"" json:"finished_at"`
	//Rating attributes
	Rating        float32
	RatingContent string
	//A message session will have a pricing
	Price float32 `gorm:"column:price"`
	//default timestamps
	// CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	// UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
