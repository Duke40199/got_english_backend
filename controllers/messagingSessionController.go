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
	"github.com/golang/got_english_backend/utils"
	"github.com/gorilla/mux"
)

func GetMessagingSessionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		//params = mux.Vars(r)
		message = "OK"
	)
	var result *[]models.MessagingSession
	var learnerID uint = 0
	var expertID uint = 0
	var err error
	//expertID
	paramsExpertID := r.URL.Query()["expert_id"]
	if len(paramsExpertID) > 0 {
		tmp, err := strconv.ParseUint(fmt.Sprint(paramsExpertID[0]), 10, 0)
		if err != nil {
			http.Error(w, "Invalid expert id.", http.StatusBadRequest)
			return
		}
		expertID = uint(tmp)
	} else {
		tmp, _ := strconv.ParseUint(fmt.Sprint(r.Context().Value("expert_id")), 10, 0)
		expertID = uint(tmp)
	}
	//learnerID
	paramsLearnerID := r.URL.Query()["learner_id"]
	if len(paramsLearnerID) > 0 {
		tmp, err := strconv.ParseUint(fmt.Sprint(paramsLearnerID[0]), 10, 0)
		if err != nil {
			http.Error(w, "Invalid learner id.", http.StatusBadRequest)
			return
		}
		learnerID = uint(tmp)
	} else {
		tmp, _ := strconv.ParseUint(fmt.Sprint(r.Context().Value("expert_id")), 10, 0)
		learnerID = uint(tmp)
	}
	messagingSession := models.MessagingSession{
		ExpertID:  &expertID,
		LearnerID: learnerID,
	}
	messagingSessionDAO := daos.GetMessagingSessionDAO()
	result, err = messagingSessionDAO.GetMessagingSessions(messagingSession)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}
func CreateMessagingSessionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	//Get learnerID
	learnerID, _ := strconv.ParseInt(fmt.Sprint(r.Context().Value("learner_id")), 10, 32)
	availableCoinCount, _ := strconv.ParseInt(fmt.Sprint(r.Context().Value("available_coin_count")), 10, 32)
	//Get messaging sessions
	messagingSession := models.MessagingSession{}
	if err := json.NewDecoder(r.Body).Decode(&messagingSession); err != nil {
		http.Error(w, "Malformed data", http.StatusBadRequest)
		return
	}
	if messagingSession.ID == "" {
		http.Error(w, "Missing (document) id.", http.StatusBadRequest)
		return
	}
	if messagingSession.PricingID == 0 {
		http.Error(w, "Missing pricing id.", http.StatusBadRequest)
		return
	}
	//Get pricing
	pricingDAO := daos.GetPricingDAO()
	pricing, _ := pricingDAO.GetPricingByID(messagingSession.PricingID)
	if availableCoinCount < int64(pricing.Price) {
		http.Error(w, "Insufficient coin.", http.StatusBadRequest)
		return
	}
	//Get exchange rate
	exchangeRateDAO := daos.GetExchangeRateDAO()
	exchangeRate, _ := exchangeRateDAO.GetExchangeRateByServiceName(config.GetServiceConfig().MessagingService)
	//Set data on model
	messagingSession.LearnerID = uint(learnerID)
	messagingSession.Pricing = *pricing
	messagingSession.ExchangeRate = *exchangeRate
	//Create
	messagingSessionDAO := daos.GetMessagingSessionDAO()
	result, err := messagingSessionDAO.CreateMessagingSession(messagingSession)
	//reduce learner available coin
	learnerDAO := daos.GetLearnerDAO()
	learner := models.Learner{
		AvailableCoinCount: uint(availableCoinCount) - pricing.Price,
	}
	_, _ = learnerDAO.UpdateLearnerByLearnerID(uint(learnerID), learner)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, result)

}
func UpdateMessagingSessionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		params  = mux.Vars(r)
		message = "OK"
	)
	//parse accountID
	messagingSession := models.MessagingSession{}
	learnerID, _ := strconv.ParseUint(fmt.Sprint(r.Context().Value("learner_id")), 10, 0)
	messagingSession.LearnerID = uint(learnerID)
	//parse body
	messagingSessionID := params["messaging_session_id"]
	if err := json.NewDecoder(r.Body).Decode(&messagingSession); err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	//Check if user inputs sessionID
	if messagingSession.ID == "" {
		http.Error(w, "missing session id.", http.StatusBadRequest)
		return
	}
	if messagingSession.IsFinished {
		http.Error(w, "Cannot update finish status using this call.", http.StatusBadRequest)
		return
	}
	//Update
	messagingSessionDAO := daos.GetMessagingSessionDAO()
	result, _, err := messagingSessionDAO.UpdateMessagingSessionByID(messagingSessionID, messagingSession)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}
func FinishMessagingSessionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		params  = mux.Vars(r)
		message = "OK"
		timeNow = time.Now()
	)
	//parse messagingSessionID
	messagingSessionID := params["messaging_session_id"]
	//Check if user inputs sessionID
	if messagingSessionID == "" {
		http.Error(w, "missing session id.", http.StatusBadRequest)
		return
	}
	//Update session isFinished
	messagingSessionDAO := daos.GetMessagingSessionDAO()
	//Check if the session is already finished
	messagingSession, _ := messagingSessionDAO.GetMessagingSessionByID(messagingSessionID)
	if messagingSession.IsFinished {
		http.Error(w, "Session is already finished.", http.StatusBadRequest)
		return
	}
	if messagingSession.IsCancelled {
		http.Error(w, "Session is already cancelled.", http.StatusBadRequest)
		return
	}
	//Update finish status
	messagingSession.IsFinished = true
	messagingSession.FinishedAt = &timeNow
	result, messagingSession, _ := messagingSessionDAO.UpdateMessagingSessionByID(messagingSessionID, *messagingSession)
	//Get coin value in VND
	pricingDAO := daos.GetPricingDAO()
	coinValue, _ := pricingDAO.GetPricings("coin_value", 0)
	coinValueInVND := (*coinValue)[0].Price
	//Calculate earning
	expertEarnings := utils.CalculateExpertEarningBySession(messagingSession.ExchangeRate.Rate, coinValueInVND, messagingSession.PaidCoins)
	earningDAO := daos.GetEarningDAO()
	earning := models.Earning{
		Value:              expertEarnings,
		ExpertID:           *messagingSession.ExpertID,
		MessagingSessionID: &messagingSession.ID,
	}
	_, err := earningDAO.CreateEarning(earning)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}
func CancelMessagingSessionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		params  = mux.Vars(r)
		message = "OK"
	)
	//parse accountID
	messagingSessionID := params["messaging_session_id"]
	messagingSession := models.MessagingSession{}
	learnerID, _ := strconv.ParseUint(fmt.Sprint(r.Context().Value("learner_id")), 10, 0)
	//parse body
	//Check if user inputs sessionID
	if messagingSession.ID != "" {
		http.Error(w, "missing session id.", http.StatusBadRequest)
		return
	} else {
		messagingSession.IsCancelled = true
	}
	//Update
	messagingSessionDAO := daos.GetMessagingSessionDAO()
	//Check if the session is already cancelled or existed
	tmpSession, _ := messagingSessionDAO.GetMessagingSessionByID(messagingSessionID)
	if tmpSession.ID == "" {
		http.Error(w, "session not found.", http.StatusBadRequest)
		return
	}
	if tmpSession.IsCancelled {
		http.Error(w, "session is already cancelled.", http.StatusBadRequest)
		return
	}
	if tmpSession.Expert != nil {
		http.Error(w, "Expert already joined this session.", http.StatusBadRequest)
		return
	}
	_, _, err := messagingSessionDAO.UpdateMessagingSessionByID(messagingSessionID, messagingSession)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	//Return the coin for learner
	pricingDAO := daos.GetPricingDAO()
	messagingPricing, _ := pricingDAO.GetPricingByID(tmpSession.PricingID)
	learnerAvailableCoin, _ := strconv.ParseUint(fmt.Sprint(r.Context().Value("available_coin_count")), 10, 32)
	currentLearner := models.Learner{
		AvailableCoinCount: uint(learnerAvailableCoin) + messagingPricing.Price,
	}
	learnerDAO := daos.GetLearnerDAO()
	result, err := learnerDAO.UpdateLearnerByLearnerID(uint(learnerID), currentLearner)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, result)

}
