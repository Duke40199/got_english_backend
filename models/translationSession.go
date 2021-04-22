package models

import (
	"time"
)

// TranslationSession model struct
type TranslationSession struct {
	ID                string `gorm:"column:id;size:255;not null; unique; primaryKey;" json:"id"`
	DurationInSeconds *uint  `gorm:"column:duration_in_seconds" json:"duration_in_seconds"`
	//TranslationSession timestamps
	IsCancelled bool       `gorm:"default:false;" json:"is_cancelled"`
	IsFinished  bool       `gorm:"default:false;" json:"is_finished"`
	StartedAt   *time.Time `gorm:"column:started_at" json:"started_at"`
	FinishedAt  *time.Time `gorm:"column:finished_at" json:"finished_at"`
	//Coins paid at the FinishedAt time
	PaidCoins uint `gorm:"column:paid_coins" json:"paid_coins"`
	//
	CreatorLearnerID uint `gorm:"column:creator_learner_id; not null" json:"creator_learner_id,omitempty"`
	//A translation session can have many learners
	Learners   []Learner `gorm:"many2many:translation_session_learners;" json:"learners"`
	LearnerIDs []uint    `gorm:"-" json:"learner_ids,omitempty"`
	//A translation session can have only one expert.
	Expert   *Expert `gorm:"foreignKey:ExpertID" json:"expert,omitempty"`
	ExpertID *uint   `gorm:"size:255" json:"expert_id,omitempty"`
	//An messaging session can only have one pricing.
	Pricing   *Pricing `gorm:"foreignKey:PricingID" json:"pricing,omitempty"`
	PricingID uint     `gorm:"size:255" json:"pricing_id"`
	//An messaging session can only have one exchange rate.
	ExchangeRate   ExchangeRate `gorm:"foreignKey:ExchangeRateID" json:"-"`
	ExchangeRateID uint         `gorm:"size:255" json:"exchange_rate_id,omitempty"`
	//An messaging session can only be rated once.
	Rating   *Rating `gorm:"foreignKey:RatingID" json:"rating,omitempty"`
	RatingID *uint   `gorm:"column:rating_id;size:255" json:"rating_id,omitempty"`
	//default timestamps
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
}
