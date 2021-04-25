package daos

import (
	"github.com/golang/got_english_backend/database"
	models "github.com/golang/got_english_backend/models"
)

type RatingAlgorithmDAO struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	TableName string
}

func (dao *RatingAlgorithmDAO) GetRatingAlgorithm() (*models.RatingAlgorithm, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	ratingAlgorithm := models.RatingAlgorithm{}
	err = db.Debug().First(&ratingAlgorithm, "id=?", 1).Error
	return &ratingAlgorithm, err
}

func (u *RatingAlgorithmDAO) UpdateRatingAlgorithm(id uint, updateInfo models.RatingAlgorithm) (int64, error) {
	db, err := database.ConnectToDB()

	if err != nil {
		return db.RowsAffected, err
	}
	result := db.Model(&models.RatingAlgorithm{}).Where("id = ?", id).
		Updates(updateInfo)
	return result.RowsAffected, result.Error
}
