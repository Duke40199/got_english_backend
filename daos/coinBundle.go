package daos

import (
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

func (u *CoinBundleDAO) CreateCoinBundle(coinBundle models.CoinBundle) (*models.CoinBundle, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Create(&coinBundle).Error
	return &coinBundle, err

}
func (u *CoinBundleDAO) GetCoinBundles() (*[]models.CoinBundle, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	coinBundles := []models.CoinBundle{}
	err = db.Debug().Model(&models.CoinBundle{}).Select("coin_bundles.*").Scan(&coinBundles).Error
	return &coinBundles, err

}
