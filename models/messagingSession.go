package models

import (
	"time"
)

// MessagingSession model struct
type MessagingSession struct {
	ID string `gorm:"column:id;size:255;not null; unique; primaryKey;" json:"id"`
	//Coins paid at the FinishedAt time
	PaidCoins uint `gorm:"column:paid_coins" json:"paid_coins"`
	//A leaner can have many message sessions.
	Learner   *Learner `gorm:"foreignKey:LearnerID" json:"learner_info,omitempty"`
	LearnerID uint     `gorm:"size:255" json:"learner_id"`
	//An expert can have many message sessions.
	Expert   *Expert `gorm:"foreignKey:ExpertID" json:"expert_info,omitempty"`
	ExpertID *uint   `gorm:"size:255" json:"expert_id,omitempty"`
	//An messaging session can only have one pricing.
	Pricing   *Pricing `gorm:"foreignKey:PricingID" json:"pricing,omitempty"`
	PricingID *uint    `gorm:"size:255" json:"pricing_id,omitempty"`
	//An messaging session can only have one exchange rate.
	ExchangeRate   ExchangeRate `gorm:"foreignKey:ExchangeRateID" json:"-"`
	ExchangeRateID uint         `gorm:"size:255" json:"exchange_rate_id,omitempty"`
	//An messaging session can only be rated once.
	Rating   *Rating `gorm:"foreignKey:RatingID" json:"rating"`
	RatingID *uint   `gorm:"size:255" json:"rating_id,omitempty"`
	//MessageSession status
	IsCancelled bool       `gorm:"default:false;" json:"is_cancelled"`
	IsFinished  bool       `gorm:"default:false;" json:"is_finished"`
	StartedAt   *time.Time `gorm:"column:started_at" json:"started_at"`
	FinishedAt  *time.Time `gorm:"column:finished_at" json:"finished_at"`

	//default timestamps
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
}
