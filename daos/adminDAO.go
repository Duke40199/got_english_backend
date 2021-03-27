package daos

import (
	"github.com/golang/got_english_backend/database"
	models "github.com/golang/got_english_backend/models"
	"github.com/google/uuid"
)

type AdminDAO struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	TableName string
}

func (dao *AdminDAO) CreateAdmin(admin models.Admin) (*models.Admin, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Create(&admin).Error
	return &admin, err

}
func (dao *AdminDAO) GetAdminByAccountID(accountID uuid.UUID) (*models.Admin, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	result := models.Admin{}
	err = db.Debug().First(&result, "account_id = ?", accountID).Error
	return &result, err
}

// func (dao *CoinBundleDAO) GetCoinBundles() (*[]models.CoinBundle, error) {
// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		return nil, err
// 	}
// 	coinBundles := []models.CoinBundle{}
// 	err = db.Debug().Model(&models.CoinBundle{}).Select("coin_bundles.*").Scan(&coinBundles).Error
// 	return &coinBundles, err

// }
// func (dao *CoinBundleDAO) UpdateCoinBundleByID(coinBundle models.CoinBundle) error {
// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		return err
// 	}
// 	err = db.Debug().Model(&coinBundle).Updates(&coinBundle).Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
