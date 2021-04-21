package models

import (
	"time"

	"github.com/google/uuid"
)

// Invoice model struct
type Invoice struct {
	ID uuid.UUID `gorm:"size:255;column:id;not null;unique; primaryKey;" json:"id"`
	//A leaner can have many invoices.
	Learner   Learner `gorm:"foreignKey:LearnerID" json:"-"`
	LearnerID uint    `gorm:"size:255" json:"learner_id"`
	//An invoice can only contain 1 coin bundle.
	CoinBundle   CoinBundle `gorm:"foreignKey:CoinBundleID" json:"-"`
	CoinBundleID uint       `gorm:"size:255" json:"coin_bundle_id"`

	PaymentMethod *string `gorm:"column:payment_method" json:"payment_method,omitempty"`
	//default timestamps
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
}
