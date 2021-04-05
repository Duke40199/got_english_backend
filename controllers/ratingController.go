package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang/got_english_backend/config"
	responseConfig "github.com/golang/got_english_backend/config"
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
	case config.GetSerivceConfig().MessagingService:
		{
			//validate whether messaging session exists.
			messagingSessionDAO := daos.GetMessagingSessionDAO()
			messagingSession, err := messagingSessionDAO.GetMessagingSessionByID(serviceID)
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
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
	case config.GetSerivceConfig().PrivateCallService:
		{
			//validate whether messaging session exists.
			privateCallDAO := daos.GetPrivateCallSessionDAO()
			privateCallSession, err := privateCallDAO.GetPrivateCallSessionByID(serviceID)
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
				return
			}
			ratingDAO := daos.GetRatingDAO()
			result, err = ratingDAO.CreatePrivateCallSessionRating(*privateCallSession, rating)
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
				return
			}
			break
		}
	case config.GetSerivceConfig().TranslationService:
		{
			//validate whether messaging session exists.
			translationSessionDAO := daos.GetTranslationSessionDAO()
			translationSession, err := translationSessionDAO.GetTranslationSessionByID(serviceID)
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
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
	responseConfig.ResponseWithSuccess(w, message, result)

}
