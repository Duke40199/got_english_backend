package daos

import (
	"fmt"

	"github.com/golang/got_english_backend/database"
	models "github.com/golang/got_english_backend/models"
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
	err = db.Debug().Create(&applicationForm).Preload("Account").Error
	return &applicationForm, err

}
func (dao *ApplicationFormDAO) GetApplicationForms(applicationForm models.ApplicationForm) (*[]models.ApplicationForm, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	query := "SELECT * from application_forms WHERE deleted_at IS NULL"
	//search by status
	if applicationForm.Status != "" {
		query += " AND status LIKE " + "'%" + applicationForm.Status + "%'"
	}
	if applicationForm.ExpertID != 0 {
		query += " AND expert_id = " + fmt.Sprint(applicationForm.ExpertID)
	}
	applicationForms := []models.ApplicationForm{}
	err = db.Debug().Model(&models.ApplicationForm{}).
		Preload("Expert").
		Preload("Expert.Account").
		Raw(query).
		Find(&applicationForms).Error
	return &applicationForms, err
}

func (dao *ApplicationFormDAO) GetApplicationFormByID(id uint) (*models.ApplicationForm, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	applicationForm := models.ApplicationForm{}
	err = db.Debug().Model(&models.ApplicationForm{}).
		Preload("Expert").
		Preload("Expert.Account").
		Select("application_forms.*").Find(&applicationForm, "id=?", id).Error
	return &applicationForm, err

}

func (dao *ApplicationFormDAO) GetApplicationFormsByExpertID(expertID uint) (*[]models.ApplicationForm, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	applicationForms := []models.ApplicationForm{}
	err = db.Debug().Order("created_at desc").Model(&models.ApplicationForm{}).
		Select("application_forms.*").Find(&applicationForms, "expert_id=?", expertID).Error
	return &applicationForms, err
}
func (dao *ApplicationFormDAO) UpdateApplicationFormByID(id uint, applicationForm models.ApplicationForm) (int64, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return db.RowsAffected, err
	}
	applicationForm.ID = id
	result := db.Model(&applicationForm).Where("id = ?", id).Updates(&applicationForm)

	return result.RowsAffected, result.Error
}
