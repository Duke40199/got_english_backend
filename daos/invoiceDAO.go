package daos

import (
	"fmt"
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
func (dao *InvoiceDAO) GetNewInvoicesCountInTimePeriod(startDate time.Time, endDate time.Time) (*[]map[string]interface{}, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	var query = "SELECT COUNT(m.created_at) AS `count` " +
		"FROM " + "(SELECT curdate() - INTERVAL (a.a + (10 * b.a) + (100 * c.a)) DAY AS Date " +
		"FROM (SELECT 0 AS a UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS a " +
		"CROSS JOIN (SELECT 0 AS a UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS b " +
		"CROSS JOIN (SELECT 0 AS a UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS c " +
		")a " + "LEFT OUTER JOIN `invoices` m " + " ON DATE(m.created_at) = a.Date " + "WHERE a.Date BETWEEN ? AND ? " +
		"GROUP BY a.Date " + "ORDER BY a.Date ASC "
	var value []int32
	err = db.Debug().
		Raw(query, startDate, endDate).
		Find(&value).
		Error
	result := make([]map[string]interface{}, len(value))
	for i := 0; i < len(value); i++ {
		result[i] = map[string]interface{}{
			fmt.Sprint(len(value)-i) + "_day_ago": value[i],
		}
	}
	return &result, err
}
