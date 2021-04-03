package daos

import (
	"time"

	"github.com/golang/got_english_backend/database"
	models "github.com/golang/got_english_backend/models"
)

type PrivateCallSessionDAO struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	TableName string
}

func (dao *PrivateCallSessionDAO) CreatePrivateCallSession(privateCallSession models.PrivateCallSession) (*models.PrivateCallSession, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Create(&privateCallSession).Error
	return &privateCallSession, err
}

func (dao *PrivateCallSessionDAO) GetCreatedPrivateCallSessionsInTimePeriod(startDate time.Time, endDate time.Time) (uint, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return 0, err
	}
	result := []models.PrivateCallSession{}
	err = db.Debug().Model(&models.PrivateCallSession{}).
		Find(&result, "created_at BETWEEN ? AND ?", startDate, endDate).Error
	return uint(len(result)), err
}

func (u *PrivateCallSessionDAO) UpdatePrivateCallSessionByID(privateCallSession models.PrivateCallSession) (int64, error) {
	db, err := database.ConnectToDB()

	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.PrivateCallSession{}).Where("id = ?", privateCallSession.ID).
		Updates(&privateCallSession)
	return result.RowsAffected, result.Error
}
