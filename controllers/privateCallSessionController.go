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

func GetPrivateCallSessionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		//params = mux.Vars(r)
		message = "OK"
	)
	var result *[]models.PrivateCallSession
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
	privateCallSession := models.PrivateCallSession{
		ExpertID:  &expertID,
		LearnerID: learnerID,
	}
	privateCallSessionDAO := daos.GetPrivateCallSessionDAO()
	result, err = privateCallSessionDAO.GetPrivateCallSessions(privateCallSession)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}
func CreatePrivateCallSessionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	//Get learnerID
	availableCoinCount, _ := strconv.ParseInt(fmt.Sprint(r.Context().Value("available_coin_count")), 10, 32)
	//Get pricing
	pricingDAO := daos.GetPricingDAO()
	pricing, _ := pricingDAO.GetPricingByID(config.GetPricingIDConfig().PrivateCallSessionPricingID)
	if availableCoinCount < int64(pricing.Price) {
		http.Error(w, "Insufficient coin.", http.StatusBadRequest)
		return
	}
	learnerID, _ := strconv.ParseInt(fmt.Sprint(r.Context().Value("learner_id")), 10, 32)
	//Get messaging sessions
	privateCallSession := models.PrivateCallSession{}
	if err := json.NewDecoder(r.Body).Decode(&privateCallSession); err != nil {
		http.Error(w, "Malformed data", http.StatusBadRequest)
		return
	}
	if privateCallSession.ID == "" {
		http.Error(w, "Missing (document) id.", http.StatusBadRequest)
		return
	}
	privateCallSession.Learner.ID = uint(learnerID)
	privateCallSession.Pricing = pricing
	//Create
	privateCallSessionDAO := daos.GetPrivateCallSessionDAO()
	result, err := privateCallSessionDAO.CreatePrivateCallSession(privateCallSession)

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, result)

}
func UpdatePrivateCallSessionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		params  = mux.Vars(r)
		message = "OK"
	)
	//parse accountID
	learnerID, _ := strconv.ParseUint(fmt.Sprint(r.Context().Value("learner_id")), 10, 0)
	privateCallSession := models.PrivateCallSession{}
	//parse body
	privateCallSessionID := params["private_call_session_id"]
	if err := json.NewDecoder(r.Body).Decode(&privateCallSession); err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	privateCallSessionDAO := daos.GetPrivateCallSessionDAO()
	result, err := privateCallSessionDAO.UpdatePrivateCallSessionByID(privateCallSessionID, privateCallSession)
	//If the session is finished, reducde learner's coin amount.
	if privateCallSession.IsFinished {
		pricingDAO := daos.GetPricingDAO()
		pricing, _ := pricingDAO.GetPricingByID(config.GetPricingIDConfig().PrivateCallSessionPricingID)
		learnerAvailableCoin, _ := strconv.ParseUint(fmt.Sprint(r.Context().Value("available_coin_count")), 10, 0)
		learner := models.Learner{
			ID:                 uint(learnerID),
			AvailableCoinCount: uint(learnerAvailableCoin) - uint(pricing.Price),
		}
		learnerDAO := daos.GetLearnerDAO()
		_, _ = learnerDAO.UpdateLearnerByLearnerID(uint(learnerID), learner)
	}
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, result)

}

func CancelPrivateCallHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		params  = mux.Vars(r)
		message = "OK"
	)
	//parse accountID

	learnerID, _ := strconv.ParseUint(fmt.Sprint(r.Context().Value("learner_id")), 10, 0)
	privateCallSessionID := params["private_call_session_id"]
	if privateCallSessionID == "" {
		http.Error(w, "missing session id.", http.StatusBadRequest)
		return
	}

	privateCallSession := models.PrivateCallSession{
		LearnerID:   uint(learnerID),
		ID:          privateCallSessionID,
		IsCancelled: true,
	}
	privateCallSessionDAO := daos.GetPrivateCallSessionDAO()
	//Check if the session is already cancelled or existed
	tmpSession, _ := privateCallSessionDAO.GetPrivateCallSessionByID(privateCallSessionID)
	if tmpSession.ID == "" {
		http.Error(w, "session not found.", http.StatusBadRequest)
		return
	}
	if tmpSession.IsCancelled {
		http.Error(w, "session is already cancelled.", http.StatusBadRequest)
		return
	}
	result, err := privateCallSessionDAO.UpdatePrivateCallSessionByID(privateCallSessionID, privateCallSession)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	config.ResponseWithSuccess(w, message, result)

}
