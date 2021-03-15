package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Message model struct
type Message struct {
	gorm.Model
	ID       uuid.UUID `gorm:"column:id; unique; primaryKey;" json:"id"`
	PhotoURL string    `json:"photo_url"`
	Content  string    `json:"content"`
	IDFrom   uuid.UUID `json:"id_from"`
	IDTo     uuid.UUID `json:"id_to"`
	//A conversation can have many messages
	MessagingSession   MessagingSession `gorm:"foreignKey:MessagingSessionID"`
	MessagingSessionID uint             `json:"messaging_session_id"`
	//default timestamps
	CreatedAt time.Time ` json:"created_at"`
	UpdatedAt time.Time ` json:"updated_at"`
}
