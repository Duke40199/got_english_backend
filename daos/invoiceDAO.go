package daos

import (
	"time"

	"github.com/golang/got_english_backend/database"
	models "github.com/golang/got_english_backend/models"
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
func (dao *InvoiceDAO) GetCreatedInvoiceInTimePeriod(startDate time.Time, endDate time.Time) (uint, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return 0, err
	}
	result := []models.Invoice{}
	err = db.Debug().Model(&models.Invoice{}).
		Find(&result, "created_at BETWEEN ? AND ?", startDate, endDate).Error
	return uint(len(result)), err
}
