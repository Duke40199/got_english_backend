package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

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
	accountID, _ := uuid.Parse(fmt.Sprint(r.Context().Value("id")))
	privateCallSession := models.PrivateCallSession{}
	//parse body
	privateCallSessionID := params["private_call_session_id"]
	if err := json.NewDecoder(r.Body).Decode(&privateCallSession); err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	privateCallSession.ID = privateCallSessionID
	learnerDAO := daos.GetLearnerDAO()
	learner, _ := learnerDAO.GetLearnerInfoByAccountID(accountID)

	privateCallSessionDAO := daos.GetPrivateCallSessionDAO()
	privateCallSession.LearnerID = learner.ID

	result, err := privateCallSessionDAO.UpdatePrivateCallSessionByID(privateCallSession)

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, result)

}

// func CancelPrivateCallHandler(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()
// 	var (
// 		params  = mux.Vars(r)
// 		message = "OK"
// 	)
// 	//parse accountID
// 	privateCallSessionID := params["private_call_id"]
// 	privateCallSession := models.PrivateCallSession{}
// 	learnerID, _ := strconv.ParseUint(fmt.Sprint(r.Context().Value("learner_id")), 10, 0)
// 	privateCallSession.LearnerID = uint(learnerID)
// 	//parse body
// 	//Check if user inputs sessionID
// 	if privateCallSession.ID != "" {
// 		http.Error(w, "missing session id.", http.StatusBadRequest)
// 		return
// 	} else {
// 		privateCallSession.ID = privateCallSessionID
// 		privateCallSession.IsCancelled = true
// 	}
// 	//Update
// 	messagingSessionDAO := daos.GetMessagingSessionDAO()
// 	//Check if the session is already cancelled or existed
// 	tmpSession, _ := messagingSessionDAO.GetMessagingSessionByID(messagingSessionID)
// 	if tmpSession.ID == "" {
// 		http.Error(w, "session not found.", http.StatusBadRequest)
// 		return
// 	}
// 	if tmpSession.IsCancelled {
// 		http.Error(w, "session is already cancelled.", http.StatusBadRequest)
// 		return
// 	}
// 	_, err := messagingSessionDAO.UpdateMessagingSessionByID(messagingSession)
// 	if err != nil {
// 		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
// 		return
// 	}
// 	//Return the coin for learner
// 	pricingDAO := daos.GetPricingDAO()
// 	messagingPricing, _ := pricingDAO.GetPricingByID(config.GetPricingIDConfig().MessagingSessionPricingID)
// 	learnerAvailableCoin, _ := strconv.ParseUint(fmt.Sprint(r.Context().Value("available_coin_count")), 10, 0)
// 	currentLearner := models.Learner{
// 		ID:                 uint(learnerID),
// 		AvailableCoinCount: uint(learnerAvailableCoin) + messagingPricing.Price,
// 	}
// 	learnerDAO := daos.GetLearnerDAO()
// 	result, err := learnerDAO.UpdateLearnerByLearnerID(currentLearner)
// 	if err != nil {
// 		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
// 		return
// 	}
// 	config.ResponseWithSuccess(w, message, result)

// }
