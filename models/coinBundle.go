package models

import (
	"time"

	"gorm.io/gorm"
)

// CoinBundle model struct
type CoinBundle struct {
	gorm.Model  `json:"-"`
	ID          uint   `gorm:"column:id;autoIncrement;not null; unique; primaryKey;" json:"id"`
	Title       string `gorm:"column:title;not null;" json:"title"`
	Description string `gorm:"column:description;not null;" json:"description"`
	Quantity    uint   `gorm:"column:quantity;not null;" json:"quantity"`
	Price       uint   `gorm:"column:price;not null;" json:"price"`
	PriceUnit   string `gorm:"column:price_unit;not null;" json:"price_unit"`
	//status
	IsDeleted bool `gorm:"column:is_deleted;not null;" json:"is_deleted"`
	//default timestamps
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
