package models

import (
	"time"

	"gorm.io/gorm"
)

// CoinBundle model struct
type CoinBundle struct {
	gorm.Model
	ID    uint   `gorm:"column:id;autoIncrement;not null; unique; primaryKey;" json:"id"`
	Title string `json:"title"`
	//An user can have many comments.
	Quantity uint `gorm:""`
	Price    uint `gorm:""`
	//default timestamps
	CreatedAt time.Time ` json:"created_at"`
	UpdatedAt time.Time ` json:"updated_at"`
}
