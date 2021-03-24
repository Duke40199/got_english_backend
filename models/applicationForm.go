package models

import (
	"time"

	"gorm.io/gorm"
)

// ApplicationForm model struct
type ApplicationForm struct {
	gorm.Model `json:"-"`
	ID         uint   `gorm:"column:id;autoIncrement;not null; unique; primaryKey;" json:"id"`
	PhotoURL   string `json:"photoUrl"`
	Status     string `gorm:"column:status;default:Pending;not null;" json:"status"`
	//An expert can have many applications
	Expert   Expert `gorm:"foreignKey:ExpertID" json:"-"`
	ExpertID uint
	//default timestamps
	CreatedAt time.Time ` json:"created_at"`
	UpdatedAt time.Time ` json:"updated_at"`
}
