package daos

import (
	"github.com/golang/got_english_backend/database"
	models "github.com/golang/got_english_backend/models"
)

type RatingDAO struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	TableName string
}

func (dao *RatingDAO) CreateMessagingSessionRating(messagingSession models.MessagingSession, rating models.Rating) (*models.Rating, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Create(&rating).Error
	if err != nil {
		return nil, err
	}
	err = db.Debug().Model(&messagingSession).Association("Rating").Append(&rating)
	return &rating, err
}

func (dao *RatingDAO) CreatePrivateCallSessionRating(privateCallSession models.PrivateCallSession, rating models.Rating) (*models.Rating, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Create(&rating).Error
	if err != nil {
		return nil, err
	}
	err = db.Debug().Model(&privateCallSession).Association("Rating").Append(&rating)
	return &rating, err
}

func (dao *RatingDAO) CreateTranslationSessionRating(translationSession models.TranslationSession, rating models.Rating) (*models.Rating, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Create(&rating).Error
	if err != nil {
		return nil, err
	}
	err = db.Debug().Model(&translationSession).Association("Rating").Append(&rating)
	return &rating, err
}
