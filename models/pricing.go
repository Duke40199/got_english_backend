package models

import (
	"time"

	"gorm.io/gorm"
)

// Pricing model struct
type Pricing struct {
	gorm.Model   `json:"-"`
	ID           uint   `gorm:"column:id;autoIncrement;not null; unique; primaryKey;" json:"id"`
	PricingName  string `gorm:"column:pricing_name;size:255" json:"pricing_name"`
	Quantity     uint   `gorm:"column:quantity;" json:"quantity"`
	QuantityUnit string `gorm:"column:quantity_unit;" json:"quantity_unit"`
	Price        uint   `gorm:"column:price;" json:"price"`
	PriceUnit    string `gorm:"column:price_unit;size:50" json:"price_unit"`

	//default timestamps

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
	//delete params
	IsDeleted bool       `gorm:"column:is_deleted" json:"is_deleted"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
