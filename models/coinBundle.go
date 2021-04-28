package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// CoinBundle model struct
type CoinBundle struct {
	gorm.Model  `json:"-"`
	ID          *uint   `gorm:"column:id;autoIncrement;not null; unique; primaryKey;" json:"id"`
	Title       *string `gorm:"column:title;not null;" json:"title,omitempty"`
	Description *string `gorm:"column:description;not null;" json:"description,omitempty"`
	Quantity    *uint   `gorm:"column:quantity;not null;" json:"quantity,omitempty"`
	Price       *uint   `gorm:"column:price;not null;" json:"price,omitempty"`
	PriceUnit   *string `gorm:"column:price_unit;not null;" json:"price_unit,omitempty"`
	//default timestamps
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
	//delete params
	IsDeleted bool       `gorm:"column:is_deleted;not null;" json:"is_deleted"`
	DeletedAt *time.Time `gorm:"column:deleted_at;" json:"deleted_at,omitempty"`
}

func VaildateCoinBundleInput(coinBundle CoinBundle) (bool, error) {
	if coinBundle.ID != nil {
		if *coinBundle.ID == 0 {
			return false, errors.New("coin bundle id cannot be 0")
		}
	}
	if coinBundle.Title != nil {
		if len(*coinBundle.Title) < 8 {
			return false, errors.New("title haves at least 8 characters")
		} else if len(*coinBundle.Title) > 55 {
			return false, errors.New("title cannot exceed 55 characters")
		}
	}
	if coinBundle.Description != nil {
		if len(*coinBundle.Description) < 8 {
			return false, errors.New("description haves at least 8 characters")
		} else if len(*coinBundle.Description) > 100 {
			return false, errors.New("description cannot exceed 100 characters")
		}
	}
	if coinBundle.Quantity != nil {
		if *coinBundle.Quantity < 1 || *coinBundle.Quantity > 1000 {
			return false, errors.New("coin quantity is from 1 to 1000")
		}
	}
	if coinBundle.Price != nil {
		return false, errors.New("you cannot edit the price")
	}
	return true, nil
}
