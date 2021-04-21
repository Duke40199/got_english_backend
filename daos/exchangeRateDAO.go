package daos

import (
	"fmt"

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

func (dao *ExchangeRateDAO) UpdateExchangeRateByID(id uint, exchangeRate models.ExchangeRate) (int64, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return db.RowsAffected, err
	}

	result := db.Model(&models.ExchangeRate{}).Where("id = ?", id).
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
		First(&result, "service_name LIKE ?", "%"+serviceName+"%").Error
	return &result, err
}

func (dao *ExchangeRateDAO) GetExchangeRates(exchangeRateQuery models.ExchangeRate) (*[]models.ExchangeRate, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	var query string = "SELECT * FROM exchange_rates WHERE id IS NOT NULL"
	if exchangeRateQuery.ID != 0 {
		query += " AND exchange_rates.id=" + fmt.Sprint(exchangeRateQuery.ID)
	}
	if len(exchangeRateQuery.ServiceName) > 0 {
		query += " AND exchange_rates.service_name LIKE " + "'" + exchangeRateQuery.ServiceName + "'"
	}
	result := []models.ExchangeRate{}
	err = db.Debug().Model(&models.ExchangeRate{}).
		Raw(query).Find(&result).Error
	return &result, err
}
