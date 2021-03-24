package daos

import (
	"github.com/golang/got_english_backend/database"
	models "github.com/golang/got_english_backend/models"
	"github.com/google/uuid"
)

type ApplicationFormDAO struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	TableName string
}

func (dao *ApplicationFormDAO) CreateApplicationForm(applicationForm models.ApplicationForm) (*models.ApplicationForm, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Create(&applicationForm).Error
	return &applicationForm, err

}
func (dao *ApplicationFormDAO) GetApplicationForms() (*[]models.ApplicationForm, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	applicationForms := []models.ApplicationForm{}
	err = db.Debug().Model(&models.ApplicationForm{}).Select("application_forms.*").Scan(&applicationForms).Error
	return &applicationForms, err

}
func (dao *ApplicationFormDAO) UpdateExpertByAccountID(accountID uuid.UUID, expertPermissions map[string]interface{}) (int64, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.Expert{}).Where("account_id = ?", accountID).
		Updates(expertPermissions)
	return result.RowsAffected, result.Error
}
