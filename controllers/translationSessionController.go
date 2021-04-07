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
	"github.com/gorilla/mux"
)

func CreateTranslationSessionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	//Get learnerID
	accountID, _ := uuid.Parse(fmt.Sprint(r.Context().Value("id")))
	availableCoinCount, _ := strconv.ParseInt(fmt.Sprint(r.Context().Value("available_coin_count")), 10, 32)
	//Get translation sessions
	translationSession := models.TranslationSession{}
	if err := json.NewDecoder(r.Body).Decode(&translationSession); err != nil {
		http.Error(w, "Malformed data", http.StatusBadRequest)
		return
	}
	if translationSession.ID == "" {
		http.Error(w, "Missing (document) id.", http.StatusBadRequest)
		return
	}
	//Check if have pricing id
	if translationSession.PricingID == 0 {
		http.Error(w, "Invalid pricing", http.StatusBadRequest)
		return
	}

	//Get pricing
	pricingDAO := daos.GetPricingDAO()
	pricing, _ := pricingDAO.GetPricingByID(translationSession.PricingID)
	if availableCoinCount < int64(pricing.Price) {
		http.Error(w, "Insufficient coin.", http.StatusBadRequest)
		return
	}
	//add creator learner
	learnerDAO := daos.GetLearnerDAO()
	learner, _ := learnerDAO.GetLearnerInfoByAccountID(accountID)
	translationSession.Learners = append(translationSession.Learners, *learner)
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
	//reduce learner available coin
	learner.AvailableCoinCount = uint(availableCoinCount) - pricing.Price
	_, _ = learnerDAO.UpdateLearnerByLearnerID(learner.ID, *learner)
	config.ResponseWithSuccess(w, message, result)

}

func UpdateTranslationSessionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		params  = mux.Vars(r)
		message = "OK"
	)
	translationSession := models.TranslationSession{}
	//parse body
	translationSessionID := params["translation_session_id"]
	if err := json.NewDecoder(r.Body).Decode(&translationSession); err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	//Check if user inputs sessionID
	if translationSession.ID != "" {
		http.Error(w, "missing session id.", http.StatusBadRequest)
		return
	}
	if len(translationSession.LearnerIDs) == 0 {
		http.Error(w, "missing learner ids.", http.StatusBadRequest)
		return
	}
	//Get learners on learnerIDS
	learnerDAO := daos.GetLearnerDAO()
	learners, err := learnerDAO.GetLearnerInfoByIDS(translationSession.LearnerIDs)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	translationSession.Learners = *learners
	//Update
	translationSessionDAO := daos.GetTranslationSessionDAO()
	result, err := translationSessionDAO.UpdateTranslationSessionByID(translationSessionID, translationSession, *learners)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}

func CancelTranslationSessionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		params  = mux.Vars(r)
		message = "OK"
	)
	//parse accountID

	learnerID, _ := strconv.ParseUint(fmt.Sprint(r.Context().Value("learner_id")), 10, 0)
	availableCoinCount, _ := strconv.ParseUint(fmt.Sprint(r.Context().Value("available_coin_count")), 10, 0)
	translationSessionID := params["translation_session_id"]
	if translationSessionID == "" {
		http.Error(w, "missing session id.", http.StatusBadRequest)
		return
	}

	translationSession := models.TranslationSession{
		ID:          translationSessionID,
		IsCancelled: true,
	}
	liveCallSessionDAO := daos.GetTranslationSessionDAO()
	//Check if the session is already cancelled or existed
	tmpSession, _ := liveCallSessionDAO.GetTranslationSessionByID(translationSession.ID)
	if tmpSession.ID == "" {
		http.Error(w, "session not found.", http.StatusBadRequest)
		return
	}
	if tmpSession.IsCancelled {
		http.Error(w, "session is already cancelled.", http.StatusBadRequest)
		return
	}
	result, err := liveCallSessionDAO.UpdateTranslationSessionByID(translationSessionID, translationSession, nil)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	//refund
	pricingDAO := daos.GetPricingDAO()
	pricing, _ := pricingDAO.GetPricingByID(tmpSession.PricingID)
	currentLearner := models.Learner{
		ID:                 uint(learnerID),
		AvailableCoinCount: uint(availableCoinCount) + pricing.Price,
	}
	learnerDAO := daos.GetLearnerDAO()
	_, _ = learnerDAO.UpdateLearnerByLearnerID(uint(learnerID), currentLearner)
	config.ResponseWithSuccess(w, message, result)

}
