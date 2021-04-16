package daos

import (
	"time"

	"github.com/golang/got_english_backend/database"
	models "github.com/golang/got_english_backend/models"
)

type TranslationSessionDAO struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	TableName string
}

//POST
func (dao *TranslationSessionDAO) CreateTranslationSession(translationSession models.TranslationSession) (*models.TranslationSession, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Create(&translationSession).Error
	return &translationSession, err
}

//GET
func (dao *TranslationSessionDAO) GetTranslationSessionByID(id string) (*models.TranslationSession, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	result := models.TranslationSession{}
	err = db.Debug().Model(&models.TranslationSession{}).
		Find(&result, "id = ?", id).Error
	return &result, err
}

func (dao *TranslationSessionDAO) GetCreatedTranslationSessionInTimePeriod(startDate time.Time, endDate time.Time) (uint, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return 0, err
	}
	result := []models.TranslationSession{}
	err = db.Debug().Model(&models.TranslationSession{}).
		Find(&result, "created_at BETWEEN ? AND ?", startDate, endDate).Error
	return uint(len(result)), err
}

func (u *TranslationSessionDAO) UpdateTranslationSessionByID(id string, translationSession models.TranslationSession, learners []models.Learner) (int64, *models.TranslationSession, error) {
	db, err := database.ConnectToDB()

	if err != nil {
		return db.RowsAffected, nil, err
	}
	//Update paid coins if finish session
	if translationSession.IsFinished {
		var tmp models.MessagingSession
		_ = db.Model(&models.MessagingSession{}).Where("id = ?", id).Select("pricing_id").First(&tmp)
		pricing, _ := pricingDAO.GetPricingByID(tmp.PricingID)
		translationSession.PaidCoins = pricing.Price
	}
	result := db.Debug().Model(&models.TranslationSession{}).Where("id = ?", id).
		Updates(&translationSession)

	//If add learner to translation session
	if len(learners) > 0 {
		_ = db.Debug().Preload("Learners").
			Where(&models.TranslationSession{ID: id}).
			Find(&translationSession).
			Association("Learners").
			Append(&learners)
	}
	return result.RowsAffected, &translationSession, result.Error
}
