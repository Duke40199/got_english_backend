package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Message model struct
type Message struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `gorm:"column:id; unique; primaryKey;" json:"id"`
	PhotoURL   string    `json:"photo_url"`
	Content    string    `json:"content"`
	IDFrom     uuid.UUID `json:"id_from"`
	IDTo       uuid.UUID `json:"id_to"`
	//A conversation can have many messages
	MessagingSession   MessagingSession `gorm:"foreignKey:MessagingSessionID"`
	MessagingSessionID string           `gorm:"messaging_session_id;size:255" json:"messaging_session_id"`

	//default timestamps
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
