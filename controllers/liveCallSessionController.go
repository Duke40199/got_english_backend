package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/models"
	"github.com/gorilla/mux"
)

func GetLiveCallSessionsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		//params = mux.Vars(r)
		message = "OK"
	)
	var result *[]models.LiveCallSession
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
	liveCallSession := models.LiveCallSession{
		ExpertID:  &expertID,
		LearnerID: learnerID,
	}
	liveCallSessionDAO := daos.GetLiveCallSessionDAO()
	result, err = liveCallSessionDAO.GetLiveCallSessions(liveCallSession)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}

func CreateLiveCallSessionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)

	//Get learnerID
	learnerID, _ := strconv.ParseInt(fmt.Sprint(r.Context().Value("learner_id")), 10, 32)
	availableCoinCount, _ := strconv.ParseInt(fmt.Sprint(r.Context().Value("available_coin_count")), 10, 32)
	//Get messaging sessions
	liveCallSession := models.LiveCallSession{}
	if err := json.NewDecoder(r.Body).Decode(&liveCallSession); err != nil {
		http.Error(w, "Malformed data", http.StatusBadRequest)
		return
	}
	if liveCallSession.ID == "" {
		http.Error(w, "Missing (document) id.", http.StatusBadRequest)
		return
	}
	if liveCallSession.PricingID == 0 {
		http.Error(w, "Missing pricing ID.", http.StatusBadRequest)
		return
	}
	//Get pricing
	pricingDAO := daos.GetPricingDAO()
	pricing, _ := pricingDAO.GetPricingByID(liveCallSession.PricingID)
	if availableCoinCount < int64(pricing.Price) {
		http.Error(w, "Insufficient coin.", http.StatusBadRequest)
		return
	}
	//Get exchange rate
	exchangeRateDAO := daos.GetExchangeRateDAO()
	exchangeRate, _ := exchangeRateDAO.GetExchangeRateByServiceName(config.GetServiceConfig().LiveCallService)

	liveCallSession.LearnerID = uint(learnerID)
	liveCallSession.Pricing = *pricing
	liveCallSession.ExchangeRate = *exchangeRate
	//Create
	liveCallSessionDAO := daos.GetLiveCallSessionDAO()
	result, err := liveCallSessionDAO.CreateLiveCallSession(liveCallSession)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	//reduce learner available coin
	learner := models.Learner{
		AvailableCoinCount: uint(availableCoinCount) - pricing.Price,
	}
	learnerDAO := daos.GetLearnerDAO()
	_, _ = learnerDAO.UpdateLearnerByLearnerID(uint(learnerID), learner)

	config.ResponseWithSuccess(w, message, result)

}
func UpdateLiveCallSessionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		params  = mux.Vars(r)
		message = "OK"
	)
	//parse accountID
	liveCallSession := models.LiveCallSession{}
	//parse body
	liveCallSessionID := params["live_call_session_id"]
	if err := json.NewDecoder(r.Body).Decode(&liveCallSession); err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	liveCallSessionDAO := daos.GetLiveCallSessionDAO()
	result, err := liveCallSessionDAO.UpdateLiveCallSessionByID(liveCallSessionID, liveCallSession)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, result)

}

func CancelLiveCallHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		params  = mux.Vars(r)
		message = "OK"
	)
	//parse accountID

	learnerID, _ := strconv.ParseUint(fmt.Sprint(r.Context().Value("learner_id")), 10, 0)
	availableCoinCount, _ := strconv.ParseUint(fmt.Sprint(r.Context().Value("available_coin_count")), 10, 0)
	liveCallSessionID := params["live_call_session_id"]
	if liveCallSessionID == "" {
		http.Error(w, "missing session id.", http.StatusBadRequest)
		return
	}

	liveCallSession := models.LiveCallSession{
		LearnerID:   uint(learnerID),
		ID:          liveCallSessionID,
		IsCancelled: true,
	}
	liveCallSessionDAO := daos.GetLiveCallSessionDAO()
	//Check if the session is already cancelled or existed
	tmpSession, _ := liveCallSessionDAO.GetLiveCallSessionByID(liveCallSession.ID)
	if tmpSession.ID == "" {
		http.Error(w, "session not found.", http.StatusBadRequest)
		return
	}
	if tmpSession.IsCancelled {
		http.Error(w, "session is already cancelled.", http.StatusBadRequest)
		return
	}
	result, err := liveCallSessionDAO.UpdateLiveCallSessionByID(liveCallSessionID, liveCallSession)
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
