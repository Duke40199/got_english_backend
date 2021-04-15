package models

import (
	"time"

	"gorm.io/gorm"
)

// LiveCallSession model struct
type LiveCallSession struct {
	gorm.Model        `json:"-"`
	ID                string `gorm:"column:id;size:255;not null; unique; primaryKey;" json:"id"`
	DurationInSeconds *uint  `gorm:"column:duration_in_seconds" json:"duration_in_seconds,omitempty"`
	//LiveCall status
	IsCancelled bool       `gorm:"default:false;" json:"is_cancelled"`
	IsFinished  bool       `gorm:"default:false;" json:"is_finished"`
	StartedAt   *time.Time `gorm:"column:started_at" json:"started_at"`
	FinishedAt  *time.Time `gorm:"column:finished_at" json:"finished_at"`
	//Coins paid at the FinishedAt time
	PaidCoins uint `gorm:"column:paid_coins" json:"paid_coins"`
	//A leaner can have many message sessions.
	Learner   Learner `gorm:"foreignKey:LearnerID" json:"-"`
	LearnerID uint    `gorm:"size:255" json:"learner_id,omitempty"`
	//An expert can have many message sessions.
	Expert   *Expert `gorm:"foreignKey:ExpertID" json:"-"`
	ExpertID *uint   `gorm:"size:255" json:"expert_id,omitempty"`
	//An messaging session can only have one pricing.
	Pricing   Pricing `gorm:"foreignKey:PricingID" json:"-"`
	PricingID uint    `gorm:"size:255" json:"pricing_id,omitempty"`
	//An messaging session can only have one exchange rate.
	ExchangeRate   ExchangeRate `gorm:"foreignKey:ExchangeRateID" json:"-"`
	ExchangeRateID uint         `gorm:"size:255" json:"exchange_rate_id,omitempty"`
	//Pricing in VND
	PricingInVND uint `gorm:"size:255" json:"pricing_in_vnd,omitempty"`
	//An messaging session can only be rated once.
	Rating   *Rating `gorm:"foreignKey:RatingID" json:"rating,omitempty"`
	RatingID *uint   `gorm:"size:255" json:"rating_id,omitempty"`
	//default timestamps
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
