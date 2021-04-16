package daos

import (
	"errors"

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

func (dao *CoinBundleDAO) GetCoinBundleByID(id uint) (*models.CoinBundle, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	coinBundle := models.CoinBundle{}
	err = db.Debug().First(&coinBundle, "id=?", id).Error
	return &coinBundle, err
}

func (dao *CoinBundleDAO) UpdateCoinBundleByID(id uint, coinBundle models.CoinBundle) (int64, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return db.RowsAffected, err
	}
	coinBundle.ID = id
	result := db.Model(&coinBundle).Where("id = ?", id).Updates(&coinBundle)

	return result.RowsAffected, result.Error
}
func (u *CoinBundleDAO) DeleteCoinBundleByID(id uint) (int64, error) {
	db, err := database.ConnectToDB()

	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.CoinBundle{}).Where("id = ?", id).
		Delete(&models.CoinBundle{})
	if result.RowsAffected == 0 {
		return result.RowsAffected, errors.New("coin bundle not found or already deleted")
	}
	return result.RowsAffected, result.Error
}
