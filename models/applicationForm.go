package models

import (
	"time"

	"gorm.io/gorm"
)

// ApplicationForm model struct
type ApplicationForm struct {
	gorm.Model `json:"-"`
	ID         uint      `gorm:"column:id;autoIncrement;not null; unique; primaryKey;" json:"id"`
	PhotoURL   string    `gorm:"column:photo_url;size:255" json:"photo_url"`
	Status     string    `gorm:"column:status;default:Pending;not null;" json:"status"`
	Type       string    `gorm:"column:type;size:255;not null" json:"type"`
	Types      *[]string `gorm:"-" json:"types,omitempty"`
	//An expert can have many applications
	Expert   *Expert `gorm:"foreignKey:ExpertID" json:"expert_info,omitempty"`
	ExpertID uint    `gorm:"" json:"expert_id,omitempty"`
	//default timestamps
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
	//delete params
	IsDeleted bool       `gorm:"column:is_deleted" json:"is_deleted"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
}
