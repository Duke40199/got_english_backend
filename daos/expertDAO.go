package daos

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/database"
	models "github.com/golang/got_english_backend/models"
	"github.com/google/uuid"
)

type ExpertDAO struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	TableName string
}

func (dao *ExpertDAO) CreateExpert(expert models.Expert) (*models.Expert, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Create(&expert).Error
	return &expert, err

}
func (dao *ExpertDAO) GetExpertByAccountID(accountID uuid.UUID) (*models.Expert, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	result := models.Expert{}
	err = db.Debug().First(&result, "account_id = ?", accountID).Error
	return &result, err
}
func (dao *ExpertDAO) GetExpertRowCount() (uint, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return 0, err
	}
	var rowsCount int64
	err = db.Debug().Model(&models.Expert{}).Count(&rowsCount).Error
	return uint(rowsCount), err
}
func (dao *ExpertDAO) GetExpertByID(expertID uint) (*models.Expert, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	result := models.Expert{}
	err = db.Debug().Preload("Account").First(&result, "id = ?", expertID).Error
	return &result, err
}
func (dao *ExpertDAO) GetExperts() (*[]models.Expert, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	result := []models.Expert{}
	err = db.Debug().Preload("Account").Find(&result).Error
	return &result, err
}
func (dao *ExpertDAO) GetCreatedExpertsInTimePeriod(startDate time.Time, endDate time.Time) (uint, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return 0, err
	}
	result := []models.Expert{}
	err = db.Debug().Model(&models.Expert{}).
		Find(&result, "experts.created_at BETWEEN ? AND ?", startDate, endDate).Error
	return uint(len(result)), err
}
func (dao *ExpertDAO) GetNewExpertsCountInTimePeriod(startDate time.Time, endDate time.Time) (*map[string]interface{}, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	var query = "SELECT COUNT(e.created_at) AS `expert_count` " +
		"FROM " + "(SELECT curdate() - INTERVAL (a.a + (10 * b.a) + (100 * c.a)) DAY AS Date " +
		"FROM (SELECT 0 AS a UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS a " +
		"CROSS JOIN (SELECT 0 AS a UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS b " +
		"CROSS JOIN (SELECT 0 AS a UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS c " +
		")a " + "LEFT OUTER JOIN `experts` e " + " ON DATE(e.created_at) = a.Date " + "WHERE a.Date BETWEEN ? AND ? " +
		"GROUP BY a.Date " + "ORDER BY a.Date ASC "
	var value []int32
	err = db.Debug().
		Raw(query, startDate, endDate).
		Find(&value).
		Error
	result := make(map[string]interface{}, len(value))
	for i := 0; i < len(value); i++ {
		result[fmt.Sprint(len(value)-i-1)+"_day_ago"] = value[i]
	}
	return &result, err
}

func (dao *ExpertDAO) UpdateExpertByAccountID(accountID uuid.UUID, expertPermissions map[string]interface{}) (int64, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.Expert{}).Where("account_id = ?", accountID).
		Updates(expertPermissions)
	return result.RowsAffected, result.Error
}

func (dao *ExpertDAO) UpdateWeightedRatingByExpertID(expertID uint, weightedRating float32) (int64, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.Expert{}).Where("id = ?", expertID).
		Updates(map[string]interface{}{"weighted_rating": weightedRating})
	return result.RowsAffected, result.Error
}

func (dao *ExpertDAO) UpdateExpertByExpertID(expertID uint, expertPermissions models.Expert) (int64, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.Expert{}).Where("id = ?", expertID).
		Updates(&expertPermissions)
	return result.RowsAffected, result.Error
}
func (dao *ExpertDAO) UpdateExpertWeightedRatingnByExpertID(expertID uint, weightedRating float32) (int64, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.Expert{}).Where("id = ?", expertID).
		Update("weighted_rating", weightedRating)
	return result.RowsAffected, result.Error
}
func (dao *ExpertDAO) GetTranslatorExperts() (*[]models.Expert, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	result := []models.Expert{}
	err = db.Debug().Model(&models.Expert{}).
		Preload("Account").
		// Raw("SELECT * FROM experts WHERE experts.id IN (SELECT experts.id FROM translation_sessions WHERE translation_sessions.is_finished = ? OR translation_sessions.is_cancelled = ?) OR experts.id NOT IN (SELECT translation_sessions.expert_id FROM got_english_db_local.translation_sessions WHERE translation_sessions.expert_id IS NOT NULL) AND experts.can_join_translation_session = ?;", true, true, true).
		Find(&result, "can_join_translation_session = ?", true).Error
	for i := 0; i < len(result); i++ {
		result[i].AverageRating, err = ratingDAO.GetExpertAverageRating(result[i])
	}

	return &result, err
}
func (dao *ExpertDAO) GetExpertSuggestions(serviceName string, limit uint) (*[]models.Expert, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	query := "weighted_rating <> 0 AND "
	switch serviceName {
	case config.GetServiceConfig().LiveCallService:
		{
			query += "can_join_live_call_session = TRUE"
			break
		}
	case config.GetServiceConfig().MessagingService:
		{
			query += "can_chat = TRUE"
			break
		}
	case config.GetServiceConfig().TranslationService:
		{
			query += "can_join_translation_session = TRUE"
			break
		}
	default:
		{
			return nil, errors.New("dao error: invalid service name")
		}
	}
	result := []models.Expert{}
	err = db.Debug().Model(&models.Expert{}).
		Preload("Account").
		Order("weighted_rating desc").
		Limit(int(limit)).
		Find(&result, query).Error
	for i := 0; i < len(result); i++ {
		result[i].AverageRating, err = ratingDAO.GetExpertAverageRating(result[i])
	}
	return &result, err

}
