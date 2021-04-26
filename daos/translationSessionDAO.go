package daos

import (
	"fmt"
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
func (dao *TranslationSessionDAO) GetTranslationSessionHistory(learnerID uint, startDate time.Time, endDate time.Time) (*[]models.TranslationSession, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	result := []models.TranslationSession{}
	err = db.Debug().Model(&models.TranslationSession{}).
		Preload("Expert").
		Preload("Pricing").
		Find(&result, "creator_learner_id = ? AND created_at BETWEEN ? AND ?", learnerID, startDate, endDate).Error
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

func (dao *TranslationSessionDAO) GetNewTranslationSessionsCountInTimePeriod(startDate time.Time, endDate time.Time) (*[]map[string]interface{}, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	var query = "SELECT COUNT(m.created_at) AS `count` " +
		"FROM " + "(SELECT curdate() - INTERVAL (a.a + (10 * b.a) + (100 * c.a)) DAY AS Date " +
		"FROM (SELECT 0 AS a UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS a " +
		"CROSS JOIN (SELECT 0 AS a UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS b " +
		"CROSS JOIN (SELECT 0 AS a UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS c " +
		")a " + "LEFT OUTER JOIN `translation_sessions` m " + " ON DATE(m.created_at) = a.Date " + "WHERE a.Date BETWEEN ? AND ? " +
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

func (u *TranslationSessionDAO) UpdateTranslationSessionByID(id string, translationSession models.TranslationSession, learners []models.Learner) (int64, *models.TranslationSession, error) {
	db, err := database.ConnectToDB()

	if err != nil {
		return db.RowsAffected, nil, err
	}
	//Update paid coins if finish session
	if translationSession.IsFinished {
		var tmp models.MessagingSession
		_ = db.Model(&models.MessagingSession{}).Where("id = ?", id).Select("pricing_id").First(&tmp)
		pricing, _ := pricingDAO.GetPricingByID(*tmp.PricingID)
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
