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
	err = db.Debug().Create(&coinBundle).Error
	if err != nil {
		return &coinBundle, err
	}
	return nil, nil
}
