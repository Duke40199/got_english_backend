package models

import (
	"time"

	"gorm.io/gorm"
)

// MessagingSession model struct
type MessagingSession struct {
	gorm.Model `json:"-"`
	ID         uint `gorm:"column:id;autoIncrement;not null; unique; primaryKey;" json:"id"`
	//Coins paid at the FinishedAt time
	PaidCoins uint `gorm:"column:paid_coins" json:"paid_coins"`
	//A leaner can have many message sessions.
	Learner   Learner `gorm:"foreignKey:LearnerID" json:"-"`
	LearnerID uint    `gorm:"size:255" json:"learner_id"`
	//An expert can have many message sessions.
	Expert   *Expert `gorm:"foreignKey:ExpertID" json:"-"`
	ExpertID *uint   `gorm:"size:255" json:"expert_id"`
	//An messaging session can only have one pricing.
	Pricing   *Pricing `gorm:"foreignKey:PricingID"`
	PricingID *uint    `gorm:"size:255" json:"pricing_id"`

	//MessageSession status
	IsCancelled bool       `gorm:"default:false;" json:"is_cancelled"`
	IsFinished  bool       `gorm:"default:false;" json:"is_finished"`
	StartedAt   *time.Time `gorm:"column:started_at" json:"started_at"`
	FinishedAt  *time.Time `gorm:"column:finished_at" json:"finished_at"`

	//default timestamps
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
