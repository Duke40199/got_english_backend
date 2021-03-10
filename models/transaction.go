package models

import (
	"time"

	"gorm.io/gorm"
)

// Transaction model struct
type Transaction struct {
	gorm.Model
	ID uint `gorm:"column:id;autoIncrement;not null; unique; primaryKey;" json:"id"`
	//A leaner can buy many coin bundles.
	Learner   Learner `gorm:"foreignKey:LearnerID"`
	LearnerID uint    `gorm:"size:255"`
	//A coin bundle can be bought by many learners.
	CoinBundle   CoinBundle `gorm:"foreignKey:CoinBundleID"`
	CoinBundleID uint       `gorm:"size:255"`
	//An user can have many comments.
	Quantity  uint      `gorm:""`
	Price     uint      `gorm:""`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime"`
}
