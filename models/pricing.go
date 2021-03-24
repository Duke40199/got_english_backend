package models

import (
	"time"

	"gorm.io/gorm"
)

// Pricing model struct
type Pricing struct {
	gorm.Model   `json:"-"`
	ID           uint   `gorm:"column:id;autoIncrement;not null; unique; primaryKey;" json:"id"`
	ServiceName  string `gorm:"column:service_name" json:"service_name"`
	Quantity     uint   `gorm:"column:quantity;" json:"quantity"`
	QuantityUnit string `gorm:"column:quantity_unit;" json:"quantity_unit"`
	Price        uint   `gorm:"column:price;" json:"price"`
	PriceUnit    string `gorm:"column:price_unit;" json:"price_unit"`
	//default timestamps
	CreatedAt time.Time ` json:"created_at"`
	UpdatedAt time.Time ` json:"updated_at"`
}
