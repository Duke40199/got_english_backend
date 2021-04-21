package models

import (
	"time"
)

// MessagingSession model struct
type ExchangeRate struct {
	ID          uint    `gorm:"column:id;not null; unique; primaryKey;" json:"id"`
	ServiceName string  `gorm:"column:service_name;size:50;unique" json:"service_name"`
	Rate        float32 `gorm:"column:float;not null" json:"rate"`
	//default timestamps
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
}
