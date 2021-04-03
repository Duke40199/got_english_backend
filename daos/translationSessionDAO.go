package daos

import (
	"time"

	"github.com/golang/got_english_backend/database"
	models "github.com/golang/got_english_backend/models"
)

type TranslationSessionDAO struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	TableName string
}

func (dao *TranslationSessionDAO) CreateTranslationSession(translationSession models.TranslationSession) (*models.TranslationSession, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Create(&translationSession).Error
	return &translationSession, err
}

func (dao *TranslationSessionDAO) GetCreatedTranslationSessionInTimePeriod(startDate time.Time, endDate time.Time) (uint, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return 0, err
	}
	result := []models.TranslationSession{}
	err = db.Debug().Model(&models.TranslationSession{}).
		Find(&result, "created_at BETWEEN ? AND ?", startDate, endDate).Error
	return uint(len(result)), err
}

func (u *TranslationSessionDAO) UpdateTranslationSessionByID(translationSession models.TranslationSession) (int64, error) {
	db, err := database.ConnectToDB()

	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.TranslationSession{}).Where("id = ?", translationSession.ID).
		Updates(&translationSession)
	return result.RowsAffected, result.Error
}
