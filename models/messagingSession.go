package models

import (
	"time"

	"gorm.io/gorm"
)

// MessagingSession model struct
type MessagingSession struct {
	gorm.Model
	ID uint `gorm:"column:id;autoIncrement;not null; unique; primaryKey;" json:"id"`
	//A leaner can have many message sessions.
	Learner   Learner `gorm:"foreignKey:LearnerID"`
	LearnerID uint    `gorm:"size:255"`
	//An expert can have many message sessions.
	Expert   Expert `gorm:"foreignKey:ExpertID"`
	ExpertID uint   `gorm:"size:255"`
	//An messaging session can only have one pricing.
	Pricing   Pricing `gorm:"foreignKey:PricingID"`
	PricingID uint    `gorm:"size:255"`
	//MessageSession status
	IsCancelled bool      `gorm:"default:false;" json:"is_cancelled"`
	IsFinished  bool      `gorm:"default:false;" json:"is_finished"`
	StartedAt   time.Time `gorm:"" json:"started_at"`
	FinishedAt  time.Time `gorm:"" json:"finished_at"`
	//Rating attributes
	Rating        float32
	RatingContent string
	//default timestamps
	CreatedAt time.Time `gorm:"column:CreatedAt;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:UpdatedAt;autoCreateTime" json:"updated_at"`
}
