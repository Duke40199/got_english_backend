package models

import (
	"time"

	"gorm.io/gorm"
)

// Pricing model struct
type Pricing struct {
	gorm.Model
	ID               uint   `gorm:"column:id;autoIncrement;not null; unique; primaryKey;" json:"id"`
	ServiceName      string `gorm:"column:service_name" json:"service_name"`
	CoinAmount       uint   `gorm:"column:coin_amount;" json:"price"`
	DurationInMinute uint   `gorm:"column:duration_in_minute;" json:"duration_in_minute"`
	//default timestamps
	CreatedAt time.Time ` json:"created_at"`
	UpdatedAt time.Time ` json:"updated_at"`
}
