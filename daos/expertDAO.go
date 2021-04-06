package daos

import (
	"time"

	"github.com/golang/got_english_backend/database"
	models "github.com/golang/got_english_backend/models"
	"github.com/google/uuid"
)

type ExpertDAO struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	TableName string
}

func (dao *ExpertDAO) CreateExpert(expert models.Expert) (*models.Expert, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Create(&expert).Error
	return &expert, err

}
func (dao *ExpertDAO) GetExpertByAccountID(accountID uuid.UUID) (*models.Expert, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	result := models.Expert{}
	err = db.Debug().First(&result, "account_id = ?", accountID).Error
	return &result, err
}

func (dao *ExpertDAO) GetCreatedExpertsInTimePeriod(startDate time.Time, endDate time.Time) (uint, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return 0, err
	}
	result := []models.Expert{}
	err = db.Debug().Model(&models.Expert{}).
		Find(&result, "experts.created_at BETWEEN ? AND ?", startDate, endDate).Error
	return uint(len(result)), err
}

func (dao *ExpertDAO) UpdateExpertByAccountID(accountID uuid.UUID, expertPermissions map[string]interface{}) (int64, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.Expert{}).Where("account_id = ?", accountID).
		Updates(expertPermissions)
	return result.RowsAffected, result.Error
}

// This will return experts which:
// 1. Is not having any translation sessions
// 2. All of their translation sessions are finished or cancelled
func (dao *ExpertDAO) GetAvailableExperts() (*[]models.Expert, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	result := []models.Expert{}
	err = db.Debug().Model(&models.Expert{}).
		Preload("Account").
		Raw("SELECT * FROM `experts` WHERE experts.id IN (SELECT experts.id FROM translation_sessions WHERE translation_sessions.is_finished = ? OR translation_sessions.is_cancelled = ?) OR experts.id NOT IN (SELECT translation_sessions.expert_id FROM translation_sessions);", true, true).
		Find(&result).Error

	return &result, err
}
