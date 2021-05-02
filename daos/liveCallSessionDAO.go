package daos

import (
	"fmt"
	"time"

	"github.com/golang/got_english_backend/database"
	models "github.com/golang/got_english_backend/models"
)

type LiveCallSessionDAO struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	TableName string
}

//POST
func (dao *LiveCallSessionDAO) CreateLiveCallSession(liveCallSession models.LiveCallSession) (*models.LiveCallSession, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Create(&liveCallSession).Error
	return &liveCallSession, err
}

//GET
func (dao *LiveCallSessionDAO) GetLiveCallSessions(liveCallSession models.LiveCallSession) (*[]models.LiveCallSession, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	result := []models.LiveCallSession{}
	var query string
	var queryValues []uint
	//ID
	if liveCallSession.LearnerID != 0 && *liveCallSession.ExpertID == 0 {
		query += "learner_id = ? "
		queryValues = append(queryValues, liveCallSession.LearnerID, 0)
	}
	if *liveCallSession.ExpertID != 0 && liveCallSession.LearnerID == 0 {
		query += "expert_id = ? "
		queryValues = append(queryValues, *liveCallSession.ExpertID, 0)
	}
	if *liveCallSession.ExpertID != 0 && liveCallSession.LearnerID != 0 {
		query += "learner_id = ? AND expert_id = ? "
		queryValues = append(queryValues, liveCallSession.LearnerID, *liveCallSession.ExpertID)
	}
	//if no input, get all
	if len(queryValues) == 0 {
		err = db.Debug().Model(&models.LiveCallSession{}).
			Preload("Rating").Preload("Learner").Preload("Expert").
			Order("created_at desc").
			Find(&result).Error
		return &result, err
	}
	//getAll
	err = db.Debug().Model(&models.LiveCallSession{}).
		Preload("Rating").Preload("Learner").Preload("Expert").
		Find(&result, query, queryValues[0], queryValues[1]).Error
	return &result, err
}

func (dao *LiveCallSessionDAO) GetLiveCallSessionByID(id string) (*models.LiveCallSession, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	result := models.LiveCallSession{}
	err = db.Debug().Model(&models.LiveCallSession{}).
		Find(&result, "id = ?", id).Error
	return &result, err
}

func (dao *LiveCallSessionDAO) GetCreatedLiveCallSessionsInTimePeriod(startDate time.Time, endDate time.Time) (uint, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return 0, err
	}
	result := []models.LiveCallSession{}
	err = db.Debug().Model(&models.LiveCallSession{}).
		Find(&result, "created_at BETWEEN ? AND ?", startDate, endDate).Error
	return uint(len(result)), err
}

func (dao *LiveCallSessionDAO) GetNewLiveCallSessionsCountInTimePeriod(startDate time.Time, endDate time.Time) (*map[string]interface{}, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	var query = "SELECT COUNT(m.created_at) AS `count` " +
		"FROM " + "(SELECT curdate() - INTERVAL (a.a + (10 * b.a) + (100 * c.a)) DAY AS Date " +
		"FROM (SELECT 0 AS a UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS a " +
		"CROSS JOIN (SELECT 0 AS a UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS b " +
		"CROSS JOIN (SELECT 0 AS a UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS c " +
		")a " + "LEFT OUTER JOIN `live_call_sessions` m " + " ON DATE(m.created_at) = a.Date " + "WHERE a.Date BETWEEN ? AND ? " +
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

func (dao *LiveCallSessionDAO) GetLiveCallSessionHistory(learnerID uint, startDate time.Time, endDate time.Time) (*[]models.LiveCallSession, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	result := []models.LiveCallSession{}
	err = db.Debug().Model(&models.LiveCallSession{}).
		Preload("Expert").
		Preload("Pricing").
		Find(&result, "learner_id = ? AND created_at BETWEEN ? AND ?  AND is_finished <> 0", learnerID, startDate, endDate).Error
	return &result, err
}

func (dao *LiveCallSessionDAO) GetLiveCallInProgress(learnerID uint, expertID uint) (bool, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return false, err
	}
	var query string
	if learnerID != 0 {
		query = "SELECT * FROM live_call_sessions WHERE learner_id = " + fmt.Sprint(learnerID) + " AND is_finished = false"
	}
	if expertID != 0 {
		query = "SELECT * FROM live_call_sessions WHERE expert_id = " + fmt.Sprint(expertID) + " AND is_finished = false"
	}
	result := []models.LiveCallSession{}
	err = db.Debug().
		Raw(query).Scan(&result).Error
	if len(result) > 0 {
		return true, nil
	}
	return false, err
}

func (u *LiveCallSessionDAO) UpdateLiveCallSessionByID(id string, liveCallSession models.LiveCallSession) (int64, *models.LiveCallSession, error) {
	db, err := database.ConnectToDB()

	if err != nil {
		return db.RowsAffected, nil, err
	}
	//Update paid coins if session is finished
	if liveCallSession.IsFinished {
		var tmp models.LiveCallSession
		_ = db.Model(&models.LiveCallSession{}).Where("id = ?", id).Select("pricing_id").First(&tmp)
		pricing, _ := pricingDAO.GetPricingByID(*tmp.PricingID)
		liveCallSession.PaidCoins = pricing.Price
	}
	result := db.Model(&models.LiveCallSession{}).Where("id = ?", id).
		Updates(&liveCallSession)
	return result.RowsAffected, &liveCallSession, result.Error
}
