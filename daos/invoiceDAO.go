package daos

import (
	"github.com/golang/got_english_backend/database"
	models "github.com/golang/got_english_backend/models"
	"github.com/google/uuid"
)

type InvoiceDAO struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	TableName string
}

func (dao *InvoiceDAO) CreateInvoice(invoice models.Invoice) (*models.Invoice, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Create(&invoice).Error
	return &invoice, err

}
func (dao *InvoiceDAO) GetExpertByAccountID(accountID uuid.UUID) (*models.Expert, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	result := models.Expert{}
	err = db.Debug().First(&result, "account_id = ?", accountID).Error
	return &result, err
}

func (dao *InvoiceDAO) UpdateExpertByAccountID(accountID uuid.UUID, expertPermissions map[string]interface{}) (int64, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.Expert{}).Where("account_id = ?", accountID).
		Updates(expertPermissions)
	return result.RowsAffected, result.Error
}
