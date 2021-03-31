package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Invoice model struct
type Invoice struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `gorm:"size:255;column:id;not null;unique; primaryKey;" json:"id"`
	//A leaner can have many invoices.
	Learner   Learner `gorm:"foreignKey:LearnerID"`
	LearnerID uint    `gorm:"size:255"`
	//An invoice can only contain 1 coin bundle.
	CoinBundle   CoinBundle `gorm:"foreignKey:CoinBundleID"`
	CoinBundleID uint       `gorm:"size:255"`

	//default timestamps
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
