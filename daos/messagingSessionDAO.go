package daos

import (
	"fmt"
	"time"

	"github.com/golang/got_english_backend/database"
	models "github.com/golang/got_english_backend/models"
)

type MessagingSessionDAO struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	TableName string
}

//POST
func (dao *MessagingSessionDAO) CreateMessagingSession(messagingSession models.MessagingSession) (*models.MessagingSession, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Create(&messagingSession).Error
	fmt.Print(err)
	return &messagingSession, err
}

//GET
func (dao *MessagingSessionDAO) GetMessagingSessionByID(id string) (*models.MessagingSession, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	result := models.MessagingSession{}
	err = db.Debug().Model(&models.MessagingSession{}).
		Find(&result, "id = ?", id).Error
	return &result, err
}

func (dao *MessagingSessionDAO) GetCreatedMessagingSessionsInTimePeriod(startDate time.Time, endDate time.Time) (uint, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return 0, err
	}
	result := []models.MessagingSession{}
	err = db.Debug().Model(&models.MessagingSession{}).
		Find(&result, "messaging_sessions.created_at BETWEEN ? AND ?", startDate, endDate).Error
	return uint(len(result)), err
}

//UPDATE
func (u *MessagingSessionDAO) UpdateMessagingSessionByID(messagingSession models.MessagingSession) (int64, error) {
	db, err := database.ConnectToDB()

	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.MessagingSession{}).Where("id = ?", messagingSession.ID).
		Updates(&messagingSession)
	return result.RowsAffected, result.Error
}

//UPDATE
func (u *MessagingSessionDAO) CancelMessagingSessionByID(messagingSession models.MessagingSession) (int64, error) {
	db, err := database.ConnectToDB()

	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.MessagingSession{}).Where("id = ?", messagingSession.ID).
		Updates(&messagingSession)
	return result.RowsAffected, result.Error
}
