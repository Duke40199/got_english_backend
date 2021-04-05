package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/models"
	"github.com/google/uuid"
)

func CreateTranaltionSessionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	//Get learnerID
	accountID, _ := uuid.Parse(fmt.Sprint(r.Context().Value("id")))
	availableCoinCount, _ := strconv.ParseInt(fmt.Sprint(r.Context().Value("available_coin_count")), 10, 32)
	//Get pricing
	pricingDAO := daos.GetPricingDAO()
	pricing, _ := pricingDAO.GetPricingByID(config.GetPricingIDConfig().TranslationSessionPricingID)
	if availableCoinCount < int64(pricing.Price) {
		http.Error(w, "Insufficient coin.", http.StatusBadRequest)
		return
	}

	//Get messaging sessions
	translationSession := models.TranslationSession{}
	if err := json.NewDecoder(r.Body).Decode(&translationSession); err != nil {
		http.Error(w, "Malformed data", http.StatusBadRequest)
		return
	}
	if translationSession.ID == "" {
		http.Error(w, "Missing (document) id.", http.StatusBadRequest)
		return
	}
	learnerDAO := daos.GetLearnerDAO()
	learner, _ := learnerDAO.GetLearnerInfoByAccountID(accountID)
	translationSession.Learners = append(translationSession.Learners, learner)
	translationSession.Pricing = *pricing
	translationSession.CreatedAt = time.Now()
	translationSession.UpdatedAt = time.Now()
	//Create
	translationSessionDAO := daos.GetTranslationSessionDAO()
	result, err := translationSessionDAO.CreateTranslationSession(translationSession)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, result)

}

// func UpdateTranlationSessionHandler(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()
// 	var (
// 		params  = mux.Vars(r)
// 		message = "OK"
// 	)
// 	//parse accountID

// 	accountID, _ := uuid.Parse(fmt.Sprint(r.Context().Value("id")))
// 	translationSession := models.TranslationSession{}
// 	//parse body
// 	translationSessionID, _ := strconv.ParseInt(params["translation_session_id"], 10, 0)
// 	if err := json.NewDecoder(r.Body).Decode(&translationSession); err != nil {
// 		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
// 		return
// 	}
// 	translationSession.ID = uint(translationSessionID)
// 	learnerDAO := daos.GetLearnerDAO()
// 	learner, _ := learnerDAO.GetLearnerInfoByAccountID(accountID)

// 	translationSessionDAO := daos.GetTranlsationSessionDAO()
// 	translationSession.Learners = learner.ID

// 	result, err := privateCallSessionDAO.UpdatePrivateCallSessionByID(privateCallSession)

// 	if err != nil {
// 		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
// 		return
// 	}
// 	config.ResponseWithSuccess(w, message, result)

// }
