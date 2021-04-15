package daos

import (
	"github.com/golang/got_english_backend/database"
	models "github.com/golang/got_english_backend/models"
)

type ExchangeRateDAO struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	TableName string
}

func (dao *ExchangeRateDAO) UpdateExchangeRate(exchangeRate models.ExchangeRate) (int64, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return db.RowsAffected, err
	}

	result := db.Model(&models.ExchangeRate{}).Where("id = ?", exchangeRate.ID).
		Updates(&exchangeRate)
	return result.RowsAffected, result.Error
}
func (dao *ExchangeRateDAO) GetExchangeRateByServiceName(serviceName string) (*models.ExchangeRate, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	result := models.ExchangeRate{}
	err = db.Debug().Model(&models.ExchangeRate{}).
		First(&result, "pricing_name LIKE ?", "%"+serviceName+"%").Error
	return &result, err
}
