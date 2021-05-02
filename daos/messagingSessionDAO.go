package daos

import (
	"fmt"
	"time"

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

//POST
func (dao *MessagingSessionDAO) CreateMessagingSession(messagingSession models.MessagingSession) (*models.MessagingSession, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Create(&messagingSession).Error
	fmt.Print(err)
	return &messagingSession, err
}

//GET
func (dao *MessagingSessionDAO) GetMessagingSessions(messagingSession models.MessagingSession) (*[]models.MessagingSession, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	result := []models.MessagingSession{}
	var query string = ""
	var queryValues []uint
	//ID
	if messagingSession.LearnerID != 0 && *messagingSession.ExpertID == 0 {
		query += "learner_id = ? "
		queryValues = append(queryValues, messagingSession.LearnerID, 0)
	}
	if *messagingSession.ExpertID != 0 && messagingSession.LearnerID == 0 {
		query += "expert_id = ? "
		queryValues = append(queryValues, *messagingSession.ExpertID, 0)
	}
	if *messagingSession.ExpertID != 0 && messagingSession.LearnerID != 0 {
		query += "learner_id = ? AND expert_id = ? "
		queryValues = append(queryValues, messagingSession.LearnerID, *messagingSession.ExpertID)
	}
	//if no input, get all
	if len(queryValues) == 0 {
		err = db.Debug().Model(&models.MessagingSession{}).
			Preload("Rating").Preload("Learner").Preload("Expert").
			Order("created_at desc").
			Find(&result).Error
		return &result, err
	}
	err = db.Debug().Model(&models.MessagingSession{}).
		Preload("Rating").Preload("Learner").Preload("Expert").
		Find(&result, query, queryValues[0], queryValues[1]).Error
	return &result, err
}

func (dao *MessagingSessionDAO) GetMessagingSessionHistory(learnerID uint, startDate time.Time, endDate time.Time) (*[]models.MessagingSession, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	result := []models.MessagingSession{}
	err = db.Debug().Model(&models.MessagingSession{}).
		Preload("Expert").
		Preload("Pricing").
		Find(&result, "learner_id = ? AND created_at BETWEEN ? AND ? AND is_finished <> 0 ", learnerID, startDate, endDate).Error
	return &result, err
}

func (dao *MessagingSessionDAO) GetMessagingSessionByID(id string) (*models.MessagingSession, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	result := models.MessagingSession{}
	err = db.Debug().Model(&models.MessagingSession{}).
		Preload("Expert").Preload("Rating").Preload("ExchangeRate").
		Find(&result, "id = ?", id).Error
	return &result, err
}

func (dao *MessagingSessionDAO) GetCreatedMessagingSessionsInTimePeriod(startDate time.Time, endDate time.Time) (uint, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return 0, err
	}
	result := []models.MessagingSession{}
	err = db.Debug().Model(&models.MessagingSession{}).
		Find(&result, "messaging_sessions.created_at BETWEEN ? AND ?", startDate, endDate).Error
	return uint(len(result)), err
}

func (dao *MessagingSessionDAO) GetNewMessagingSessionsCountInTimePeriod(startDate time.Time, endDate time.Time) (*map[string]interface{}, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	var query = "SELECT COUNT(m.created_at) AS `count` " +
		"FROM " + "(SELECT curdate() - INTERVAL (a.a + (10 * b.a) + (100 * c.a)) DAY AS Date " +
		"FROM (SELECT 0 AS a UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS a " +
		"CROSS JOIN (SELECT 0 AS a UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS b " +
		"CROSS JOIN (SELECT 0 AS a UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS c " +
		")a " + "LEFT OUTER JOIN `messaging_sessions` m " + " ON DATE(m.created_at) = a.Date " + "WHERE a.Date BETWEEN ? AND ? " +
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
func (dao *MessagingSessionDAO) GetMessagingInProgress(learnerID uint, expertID uint) (bool, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return false, err
	}
	var query string
	if learnerID != 0 {
		query = "SELECT * FROM messaging_sessions WHERE learner_id = " + fmt.Sprint(learnerID) + " AND is_finished = false"
	}
	if expertID != 0 {
		query = "SELECT * FROM messaging_sessions WHERE expert_id = " + fmt.Sprint(expertID) + " AND is_finished = false"
	}
	result := []models.MessagingSession{}
	err = db.Debug().
		Raw(query).Scan(&result).Error
	if len(result) > 0 {
		return true, nil
	}
	return false, err
}

//UPDATE
func (u *MessagingSessionDAO) UpdateMessagingSessionByID(id string, messagingSession models.MessagingSession) (int64, *models.MessagingSession, error) {
	db, err := database.ConnectToDB()

	if err != nil {
		return db.RowsAffected, nil, err
	}
	//Update paid coins
	if messagingSession.IsFinished {
		var tmp models.MessagingSession
		_ = db.Model(&models.MessagingSession{}).Where("id = ?", id).Select("pricing_id").First(&tmp)
		pricing, _ := pricingDAO.GetPricingByID(*tmp.PricingID)
		messagingSession.PaidCoins = pricing.Price
	}
	result := db.Model(&models.MessagingSession{}).Where("id = ?", id).
		Updates(&messagingSession)
	return result.RowsAffected, &messagingSession, result.Error
}

//UPDATE
func (u *MessagingSessionDAO) CancelMessagingSessionByID(messagingSession models.MessagingSession) (int64, error) {
	db, err := database.ConnectToDB()

	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.MessagingSession{}).Where("id = ?", messagingSession.ID).
		Updates(&messagingSession)
	return result.RowsAffected, result.Error
}
