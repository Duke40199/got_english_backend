package models

import (
	"time"

	"gorm.io/gorm"
)

// Expert model struct
type Rating struct {
	gorm.Model `json:"-"`

	ID      uint    `gorm:"column:id;not null;unique; primaryKey;" json:"id"`
	Score   float32 `gorm:"column:score" json:"score"`
	Comment string  `gorm:"size:225;column:comment;" json:"comment"`

	//Rating is created by one learner.
	Learner   Learner `gorm:"foreignKey:LearnerID;" json:"learner,omitempty"`
	LearnerID uint    `gorm:"size:225;column:learner_id" json:"learner_id"`

	//A rating can only be in a session
	LiveCallSession    *LiveCallSession    `json:"live_call_session,omitempty"`
	MessagingSession   *MessagingSession   `json:"messaging_session,omitempty"`
	TranslationSession *TranslationSession `json:"translation_session,omitempty"`

	//default timestamps
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
