package daos

import (
	"fmt"
	"strconv"
	"time"

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

func (dao *LearnerDAO) GetLearnerInfoByIDS(learnerIDS []uint) (*[]models.Learner, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	//Initial learner
	var query string = "id = " + strconv.FormatUint(uint64(learnerIDS[0]), 10)
	for i := 1; i < len(learnerIDS); i++ {
		query += " OR id = " + strconv.FormatUint(uint64(learnerIDS[i]), 10)
	}
	learners := []models.Learner{}
	err = db.Debug().Find(&learners, query).Error
	return &learners, err
}

func (dao *LearnerDAO) GetCreatedLearnersInTimePeriod(startDate time.Time, endDate time.Time) (uint, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return 0, err
	}
	result := []models.Learner{}
	err = db.Debug().Model(&models.Learner{}).
		Find(&result, "learners.created_at BETWEEN ? AND ?", startDate, endDate).Error
	return uint(len(result)), err
}

func (dao *LearnerDAO) GetNewLearnersCountInTimePeriod(startDate time.Time, endDate time.Time) (*[]map[string]interface{}, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	var query = "SELECT COUNT(l.created_at) AS `count` " +
		"FROM " + "(SELECT curdate() - INTERVAL (a.a + (10 * b.a) + (100 * c.a)) DAY AS Date " +
		"FROM (SELECT 0 AS a UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS a " +
		"CROSS JOIN (SELECT 0 AS a UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS b " +
		"CROSS JOIN (SELECT 0 AS a UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS c " +
		")a " + "LEFT OUTER JOIN `learners` l " + " ON DATE(l.created_at) = a.Date " + "WHERE a.Date BETWEEN ? AND ? " +
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

func (dao *LearnerDAO) UpdateLearnerByLearnerID(learnerID uint, learner models.Learner) (int64, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.Learner{}).Where("id = ?", learnerID).
		Updates(&learner)
	return result.RowsAffected, result.Error
}
