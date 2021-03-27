package daos

import (
	"github.com/golang/got_english_backend/database"
	models "github.com/golang/got_english_backend/models"
	"github.com/google/uuid"
)

type ModeratorDAO struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	TableName string
}

func (dao *ModeratorDAO) CreateModerator(moderator models.Moderator) (*models.Moderator, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Create(&moderator).Error
	return &moderator, err

}

func (dao *ModeratorDAO) UpdateModeratorByAccountID(accountID uuid.UUID, moderatorPermissions map[string]interface{}) (int64, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.Moderator{}).Where("account_id = ?", accountID).
		Updates(moderatorPermissions)
	return result.RowsAffected, result.Error
}
func (dao *ModeratorDAO) GetModeratorByAccountID(accountID uuid.UUID) (*models.Moderator, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	result := models.Moderator{}
	err = db.Debug().First(&result, "account_id = ?", accountID).Error
	return &result, err
}
