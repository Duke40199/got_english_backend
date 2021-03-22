package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Invoice model struct
type Invoice struct {
	gorm.Model
	ID uuid.UUID `gorm:"size:255;column:id;not null;unique; primaryKey;" json:"id"`
	//A leaner can have many invoices.
	Learner   Learner `gorm:"foreignKey:LearnerID"`
	LearnerID uint    `gorm:"size:255"`
	//An invoice can only contain 1 coin bundle.
	CoinBundle   CoinBundle `gorm:"foreignKey:CoinBundleID"`
	CoinBundleID uint       `gorm:"size:255"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime"`
}
