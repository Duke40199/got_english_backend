package daos

import (
	"time"

	"github.com/golang/got_english_backend/database"
	models "github.com/golang/got_english_backend/models"
)

type EarningDAO struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	TableName string
}

func (dao *EarningDAO) CreateEarning(earning models.Earning) (*models.Earning, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Create(&earning).Error
	return &earning, err

}
func (dao *EarningDAO) GetEarningByExpertID(expertID uint, startDate time.Time, endDate time.Time) (*[]models.Earning, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	result := []models.Earning{}
	err = db.Debug().Model(&models.Earning{}).
		Where("expert_id =?", expertID).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Find(&result).Error
	return &result, err
}
