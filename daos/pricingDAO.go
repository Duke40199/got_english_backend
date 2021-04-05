package daos

import (
	"github.com/golang/got_english_backend/database"
	models "github.com/golang/got_english_backend/models"
)

type PricingDAO struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	TableName string
}

func (dao *PricingDAO) GetPricingByID(id uint) (*models.Pricing, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	pricing := models.Pricing{}
	err = db.Debug().First(&pricing, "id=?", id).Error
	return &pricing, err
}
