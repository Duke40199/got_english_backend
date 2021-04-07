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

func (dao *AdminDAO) UpdateAdminByAccountID(accountID uuid.UUID, permissions map[string]interface{}) (int64, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.Admin{}).Where("account_id = ?", accountID).
		Updates(permissions)
	return result.RowsAffected, result.Error
}
