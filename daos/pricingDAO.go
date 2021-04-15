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

func (dao *PricingDAO) GetPricings(pricingName string, id uint) (*[]models.Pricing, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	pricing := []models.Pricing{}
	// If id is not inputted
	if id == 0 {
		err = db.Debug().Where("pricings.pricing_name LIKE ?", "%"+pricingName+"%").Find(&pricing).Error
	} else {
		err = db.Debug().Where("pricings.pricing_name LIKE ? AND pricings.id = ?", "%"+pricingName+"%", id).Find(&pricing).Error
	}
	return &pricing, err
}

func (u *PricingDAO) CreatePricingHandler(pricing models.Pricing) (*models.Pricing, error) {
	db, err := database.ConnectToDB()

	if err != nil {
		return nil, err
	}
	err = db.Debug().Model(&models.Pricing{}).Create(&pricing).Error
	return &pricing, err
}

func (u *PricingDAO) UpdatePricingByID(id uint, updateInfo models.Pricing) (int64, error) {
	db, err := database.ConnectToDB()

	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.Pricing{}).Where("id = ?", id).
		Updates(updateInfo)
	return result.RowsAffected, result.Error
}

func (u *PricingDAO) DeletePricingByID(id uint) (int64, error) {
	db, err := database.ConnectToDB()

	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.Pricing{}).Where("id = ?", id).
		Delete(&models.Pricing{})
	return result.RowsAffected, result.Error
}
