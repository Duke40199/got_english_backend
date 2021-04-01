package daos

import (
	"fmt"

	"github.com/golang/got_english_backend/database"
	models "github.com/golang/got_english_backend/models"
)

type PrivateCallSessionDAO struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	TableName string
}

func (dao *PrivateCallSessionDAO) CreatePrivateCallSession(privateCallSession models.PrivateCallSession) (*models.PrivateCallSession, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Create(&privateCallSession).Error
	fmt.Print(err)
	return &privateCallSession, err

}
func (u *PrivateCallSessionDAO) UpdatePrivateCallSessionByID(privateCallSession models.PrivateCallSession) (int64, error) {
	db, err := database.ConnectToDB()

	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.PrivateCallSession{}).Where("id = ?", privateCallSession.ID).
		Updates(&privateCallSession)
	return result.RowsAffected, result.Error
}
