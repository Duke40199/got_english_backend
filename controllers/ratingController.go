package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/models"
)

func CreateRatingHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	var result interface{}
	ctx := r.Context()
	learnerID, _ := strconv.ParseUint(fmt.Sprint(ctx.Value("learner_id")), 10, 0)
	ratingInfo := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&ratingInfo); err != nil {
		fmt.Print(err)
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	if ratingInfo["service_id"] == nil {
		http.Error(w, "Missing service id", http.StatusBadRequest)
		return
	}
	serviceID := fmt.Sprint(ratingInfo["service_id"])
	score, _ := strconv.ParseFloat(fmt.Sprint(ratingInfo["score"]), 10)

	if score < 1 || score > 5 {
		http.Error(w, "Score is between 1 and 5.", http.StatusBadRequest)
		return
	}

	rating := models.Rating{
		Score:     float32(score),
		Comment:   fmt.Sprint(ratingInfo["comment"]),
		LearnerID: uint(learnerID),
	}
	switch ratingInfo["service"] {
	case config.GetServiceConfig().MessagingService:
		{
			//validate whether messaging session exists.
			messagingSessionDAO := daos.GetMessagingSessionDAO()
			messagingSession, err := messagingSessionDAO.GetMessagingSessionByID(serviceID)
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
				return
			}
			//validate whether session is finished before rating
			if !messagingSession.IsFinished {
				http.Error(w, "Cannot rate a session which is not finished.", http.StatusBadRequest)
				return
			}
			//validate whether session is rated.
			if messagingSession.Rating != nil {
				http.Error(w, "Session is already rated.", http.StatusBadRequest)
				return
			}
			ratingDAO := daos.GetRatingDAO()
			result, err = ratingDAO.CreateMessagingSessionRating(*messagingSession, rating)
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
				return
			}
			break
		}
	case config.GetServiceConfig().LiveCallService:
		{
			//validate whether messaging session exists.
			liveCallDAO := daos.GetLiveCallSessionDAO()
			liveCallSession, err := liveCallDAO.GetLiveCallSessionByID(serviceID)
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
				return
			}
			//validate whether session is finished before rating
			if !liveCallSession.IsFinished {
				http.Error(w, "Cannot rate a session which is not finished.", http.StatusBadRequest)
				return
			}
			//validate whether session is rated.
			if liveCallSession.Rating != nil {
				http.Error(w, "Session is already rated.", http.StatusBadRequest)
				return
			}
			ratingDAO := daos.GetRatingDAO()
			result, err = ratingDAO.CreateLiveCallSessionRating(*liveCallSession, rating)
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
				return
			}
			break
		}
	case config.GetServiceConfig().TranslationService:
		{
			//validate whether messaging session exists.
			translationSessionDAO := daos.GetTranslationSessionDAO()
			translationSession, err := translationSessionDAO.GetTranslationSessionByID(serviceID)
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
				return
			}
			//validate whether session is finished before rating
			if !translationSession.IsFinished {
				http.Error(w, "Cannot rate a session which is not finished.", http.StatusBadRequest)
				return
			}
			//validate whether session is rated.
			if translationSession.Rating != nil {
				http.Error(w, "Session is already rated.", http.StatusBadRequest)
				return
			}
			ratingDAO := daos.GetRatingDAO()
			result, err = ratingDAO.CreateTranslationSessionRating(*translationSession, rating)
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
				return
			}
			break
		}
	default:
		{
			http.Error(w, "Invalid/missing service", http.StatusBadRequest)
			return
		}
	}
	config.ResponseWithSuccess(w, message, result)

}

func GetRatingsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	var expertID uint = 0
	var err error
	var result *[]models.Rating
	if len(r.URL.Query()["expert_id"]) > 0 {
		tmp, err := strconv.ParseUint(r.URL.Query()["expert_id"][0], 10, 0)
		if err != nil {
			http.Error(w, "Invalid expert id.", http.StatusBadRequest)
			return
		}
		expertID = uint(tmp)
	}

	ratingDAO := daos.GetRatingDAO()
	result, err = ratingDAO.GetRatings(expertID)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}
