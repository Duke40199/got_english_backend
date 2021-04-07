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

func (dao *RatingDAO) CreateLiveCallSessionRating(liveCallSession models.LiveCallSession, rating models.Rating) (*models.Rating, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	err = db.Debug().Create(&rating).Error
	if err != nil {
		return nil, err
	}
	err = db.Debug().Model(&liveCallSession).Association("Rating").Append(&rating)
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
func (dao *RatingDAO) GetRatings() (*[]models.Rating, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return nil, err
	}
	var result []models.Rating
	err = db.Debug().Model(&models.Rating{}).
		Preload("Learner").
		Preload("TranslationSession").
		Preload("LiveCallSession").
		Preload("MessagingSession.Expert").
		Preload("MessagingSession.Expert.Account").
		Order("created_at desc").
		Find(&result).Error
	return &result, err
}
func (dao *RatingDAO) GetExpertAverageRating(expert models.Expert) (float32, error) {
	db, err := database.ConnectToDB()
	if err != nil {
		return 0, err
	}
	var averageRating float32 = 0.0
	var sumRating float32
	var ratings []float32

	//messaging sessions
	expertMessagingSessions := []models.MessagingSession{}
	err = db.Debug().Model(&models.MessagingSession{}).Preload("Rating").
		Find(&expertMessagingSessions, "expert_id=?", expert.ID).Error
	//live call sessions
	expertLiveCallSessions := []models.LiveCallSession{}
	err = db.Debug().Model(&models.LiveCallSession{}).Preload("Rating").
		Find(&expertLiveCallSessions, "expert_id=?", expert.ID).Error
	//translation call sessions
	expertTranslationSessions := []models.TranslationSession{}
	err = db.Debug().Model(&models.TranslationSession{}).Preload("Rating").
		Find(&expertTranslationSessions, "expert_id=?", expert.ID).Error
	//combine all ratings
	if len(expertMessagingSessions) > 0 {
		for i := 0; i < len(expertMessagingSessions); i++ {
			//If the session is rated.
			if expertMessagingSessions[i].Rating != nil {
				if expertMessagingSessions[i].Rating.Score > 0 {
					ratings = append(ratings, expertMessagingSessions[i].Rating.Score)
				}
			}
		}
	}
	if len(expertLiveCallSessions) > 0 {
		for i := 0; i < len(expertLiveCallSessions); i++ {
			//If the session is rated.
			if expertLiveCallSessions[i].Rating != nil {
				if expertLiveCallSessions[i].Rating.Score > 0 {
					ratings = append(ratings, expertLiveCallSessions[i].Rating.Score)
				}
			}
		}
	}
	if len(expertTranslationSessions) > 0 {
		for i := 0; i < len(expertTranslationSessions); i++ {
			//If the session is rated.
			if expertTranslationSessions[i].Rating != nil {
				if expertTranslationSessions[i].Rating.Score > 0 {
					ratings = append(ratings, expertTranslationSessions[i].Rating.Score)
				}
			}
		}
	}
	if len(ratings) > 0 {
		for i := 0; i < len(ratings); i++ {
			sumRating += ratings[i]
		}
		averageRating = sumRating / float32(len(ratings))
	}
	return averageRating, err
}
