package daos

import (
	"errors"
	"fmt"

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
	var result = db
	pricingInfo := models.Pricing{}
	_ = db.Model(&models.Pricing{}).First(&pricingInfo, "id = ?", id)
	//check if the pricing is a coin value
	if pricingInfo.PricingName == "coin_value" {
		result = db.Model(&models.Pricing{}).Where("id = ?", id).
			Updates(map[string]interface{}{"price": updateInfo.Price})
		//update coin bundles based on new price value
		coinBundles := []models.CoinBundle{}
		var bundleUpdateQuery = "INSERT into `coin_bundles` (id, price,price_unit) VALUES "
		_ = db.Model(&models.CoinBundle{}).Find(&coinBundles)
		for i := 0; i < len(coinBundles); i++ {
			coinNewValueQuery := "(" + fmt.Sprint(i+1) + "," + fmt.Sprint(updateInfo.Price*(*coinBundles[i].Quantity)) + "," + "'VND') "
			if i < len(coinBundles)-1 {
				coinNewValueQuery += ", "
			}
			bundleUpdateQuery += coinNewValueQuery
		}
		bundleUpdateQuery += "ON DUPLICATE KEY UPDATE price = VALUES(price);"
		_ = db.Exec(bundleUpdateQuery)
	} else {
		result = db.Model(&models.Pricing{}).Where("id = ?", id).
			Updates(updateInfo)
	}

	return result.RowsAffected, result.Error
}

func (u *PricingDAO) DeletePricingByID(id uint) (int64, error) {
	db, err := database.ConnectToDB()

	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.Pricing{}).Where("id = ?", id).
		Updates(&models.Pricing{IsDeleted: true}).
		Delete(&models.Pricing{})
	if (result.RowsAffected) == 0 {
		return result.RowsAffected, errors.New("pricing not found or already deleted")
	}
	return result.RowsAffected, result.Error
}
