package models

import (
	"time"
)

// MessagingSession model struct
type Earning struct {
	ID    uint    `gorm:"column:id;not null; unique; primaryKey;" json:"id"`
	Value float32 `gorm:"column:value;not null" json:"value"`
	Unit  string  `gorm:"column:unit;default:VND" json:"unit"`
	//An earning can only belongs to one expert
	Expert   *Expert `gorm:"foreignKey:ExpertID" json:"expert_info,omitempty"`
	ExpertID uint    `gorm:"column:expert_id; not null" json:"expert_id"`
	//An earning can belong to one messagingSession
	MessagingSession   *MessagingSession `gorm:"foreignKey:MessagingSessionID" json:"messaging_session,omitempty"`
	MessagingSessionID *string           `gorm:"column:messaging_session_id;" json:"messaging_session_id,omitempty"`
	//An earning can belong to one messagingSession
	LiveCallSession   *LiveCallSession `gorm:"foreignKey:LiveCallSessionID" json:"live_call_session,omitempty"`
	LiveCallSessionID *string          `gorm:"column:live_call_session_id;" json:"live_call_session_id,omitempty"`
	//An earning can belong to one messagingSession
	TranslationSession   *TranslationSession `gorm:"foreignKey:TranslationSessionID" json:"translation_session,omitempty"`
	TranslationSessionID *string             `gorm:"column:translation_session_id;" json:"translation_session_id,omitempty"`
	//default timestamps
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}
