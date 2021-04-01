package daos

import (
	"github.com/golang/got_english_backend/database"
	models "github.com/golang/got_english_backend/models"
	"github.com/google/uuid"
)

type LearnerDAO struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	TableName string
}

func (dao *LearnerDAO) CreateLearner(learner models.Learner) (*models.Learner, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Create(&learner).Error
	return &learner, err

}

func (dao *LearnerDAO) GetLearnerInfoByAccountID(accountID uuid.UUID) (*models.Learner, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}

	learner := models.Learner{}
	err = db.Debug().First(&learner, "account_id = ?", accountID).Error
	return &learner, err

}

func (dao *LearnerDAO) UpdateLearnerByLearnerID(learner models.Learner) (int64, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.Learner{}).Where("id = ?", learner.ID).
		Updates(&learner)
	return result.RowsAffected, result.Error
}
