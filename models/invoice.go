package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Invoice model struct
type Invoice struct {
	gorm.Model
	ID uuid.UUID `gorm:"size:255;column:id;not null;unique; primaryKey;" json:"id"`
	//A leaner can have many invoices.
	Learner   Learner `gorm:"foreignKey:LearnerID"`
	LearnerID uint    `gorm:"size:255"`
	//An invoice can only contain 1 coin bundle.
	CoinBundle   CoinBundle `gorm:"foreignKey:CoinBundleID"`
	CoinBundleID uint       `gorm:"size:255"`
	//An invoice can only contain 1 messaging session.
	MessagingSession   MessagingSession `gorm:"foreignKey:MessagingSessionID"`
	MessagingSessionID uint             `gorm:"size:255"`
	//An invoice can only contain 1 translation session.
	TranslationSession   TranslationSession `gorm:"foreignKey:TranslationSessionID"`
	TranslationSessionID uint               `gorm:"size:255"`
	//An invoice can only contain 1 private call session.
	PrivateCallSession   PrivateCallSession `gorm:"foreignKey:PrivateCallSessionID"`
	PrivateCallSessionID uint               `gorm:"size:255"`
	//An user can have many comments.
	Price     uint      `gorm:"column:price"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime"`
}
