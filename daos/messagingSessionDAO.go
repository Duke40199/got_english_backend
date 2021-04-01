package daos

import (
	"fmt"

	"github.com/golang/got_english_backend/database"
	models "github.com/golang/got_english_backend/models"
)

type MessagingSessionDAO struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	TableName string
}

func (dao *MessagingSessionDAO) CreateMessagingSession(messagingSession models.MessagingSession) (*models.MessagingSession, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Create(&messagingSession).Error
	fmt.Print(err)
	return &messagingSession, err

}
func (u *MessagingSessionDAO) UpdateMessagingSessionByID(messagingSession models.MessagingSession) (int64, error) {
	db, err := database.ConnectToDB()

	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.MessagingSession{}).Where("id = ?", messagingSession.ID).
		Updates(&messagingSession)
	return result.RowsAffected, result.Error
}

// func (dao *InvoiceDAO) GetExpertByAccountID(accountID uuid.UUID) (*models.Expert, error) {
// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		return nil, err
// 	}
// 	result := models.Expert{}
// 	err = db.Debug().First(&result, "account_id = ?", accountID).Error
// 	return &result, err
// }

// func (dao *InvoiceDAO) UpdateExpertByAccountID(accountID uuid.UUID, expertPermissions map[string]interface{}) (int64, error) {
// 	db, err := database.ConnectToDB()
// 	if err != nil {
// 		return db.RowsAffected, err
// 	}
// 	result := db.Model(&models.Expert{}).Where("account_id = ?", accountID).
// 		Updates(expertPermissions)
// 	return result.RowsAffected, result.Error
// }
