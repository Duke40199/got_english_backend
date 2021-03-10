package models

import (
	"time"

	"gorm.io/gorm"
)

// Pricing model struct
type Pricing struct {
	gorm.Model
	ID          uint   `gorm:"column:id;autoIncrement;not null; unique; primaryKey;" json:"id"`
	ServiceType string `json:"service_type"`
	Price       uint   `gorm:"column:price;" json:"price"`
	//default timestamps
	CreatedAt time.Time ` json:"created_at"`
	UpdatedAt time.Time ` json:"updated_at"`
}
