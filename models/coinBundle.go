package models

import (
	"time"

	"gorm.io/gorm"
)

// CoinBundle model struct
type CoinBundle struct {
	gorm.Model
	ID          uint   `gorm:"column:id;autoIncrement;not null; unique; primaryKey;" json:"id"`
	Title       string `gorm:"column:title;" json:"title"`
	Description string `gorm:"column:description;" json:"description"`
	Quantity    uint   `gorm:"column:quantity;" json:"quantity"`
	Price       uint   `gorm:"column:price;" json:"price"`
	PriceUnit   string `gorm:"column:price_unit;" json:"price_unit"`
	//default timestamps
	CreatedAt time.Time ` json:"created_at"`
	UpdatedAt time.Time ` json:"updated_at"`
}
