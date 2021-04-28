package daos

import (
	"errors"
	"fmt"

	"github.com/golang/got_english_backend/database"
	models "github.com/golang/got_english_backend/models"
)

type CoinBundleDAO struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	TableName string
}

func (dao *CoinBundleDAO) CreateCoinBundle(coinBundle models.CoinBundle) (*models.CoinBundle, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	pricingDAO := GetPricingDAO()
	pricing, _ := pricingDAO.GetPricings("coin_value", 0)
	//Calculate pricing
	bundlePrice := (*pricing)[0].Price * *coinBundle.Quantity
	coinBundle.Price = &bundlePrice
	err = db.Debug().Create(&coinBundle).Error
	return &coinBundle, err

}
func (dao *CoinBundleDAO) GetCoinBundles(id uint) (*[]models.CoinBundle, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	coinBundles := []models.CoinBundle{}
	// if id ==0 => get all
	if id == 0 {
		err = db.Debug().Model(&models.CoinBundle{}).Select("coin_bundles.*").Scan(&coinBundles).Error
		return &coinBundles, err
	}
	err = db.Debug().Model(&models.CoinBundle{}).Select("coin_bundles.*").Where("id = ?", id).Scan(&coinBundles).Error
	return &coinBundles, err
}

func (dao *CoinBundleDAO) GetCoinBundleByID(id uint, coinBundle models.CoinBundle) (*models.CoinBundle, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	var query string = "SELECT * from coin_bundles WHERE (is_deleted IS FALSE"
	//if query for deleted entries
	if coinBundle.IsDeleted {
		query += " OR is_deleted IS TRUE)"
	} else {
		query += ")"
	}
	if id != 0 {
		query += " AND id =" + fmt.Sprint(id)
	}
	err = db.Debug().Model(&models.CoinBundle{}).Raw(query).Find(&coinBundle).Error
	return &coinBundle, err
}

func (dao *CoinBundleDAO) UpdateCoinBundleByID(id uint, coinBundle models.CoinBundle) (int64, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return db.RowsAffected, err
	}
	coinBundle.ID = &id
	pricingDAO := GetPricingDAO()
	pricing, _ := pricingDAO.GetPricings("coin_value", 0)
	//Calculate pricing
	bundlePrice := (*pricing)[0].Price * *coinBundle.Quantity
	coinBundle.Price = &bundlePrice
	result := db.Model(&coinBundle).Where("id = ?", id).Updates(&coinBundle)

	return result.RowsAffected, result.Error
}
func (u *CoinBundleDAO) DeleteCoinBundleByID(id uint) (int64, error) {
	db, err := database.ConnectToDB()

	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.CoinBundle{}).Where("id = ?", id).
		Updates(&models.CoinBundle{IsDeleted: true}).
		Delete(&models.CoinBundle{})
	if result.RowsAffected == 0 {
		return result.RowsAffected, errors.New("coin bundle not found or already deleted")
	}
	return result.RowsAffected, result.Error
}
