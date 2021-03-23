package daos

type InvoiceDAO struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	TableName string
}

// func (dao *CoinBundleDAO) CreateCoinBundle(coinBundle models.CoinBundle) (*models.CoinBundle, error) {
// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = db.Debug().Create(&coinBundle).Error
// 	return &coinBundle, err

// }
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
