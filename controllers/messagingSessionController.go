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

func CreateMessagingSessionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	//Get learnerID
	availableCoinCount, _ := strconv.ParseInt(fmt.Sprint(r.Context().Value("available_coin_count")), 10, 32)
	//Get pricing
	pricingDAO := daos.GetPricingDAO()
	pricing, _ := pricingDAO.GetPricingByID(config.GetPricingIDConfig().MessagingSessionPricingID)
	if availableCoinCount < int64(pricing.Price) {
		http.Error(w, "Insufficient coin.", http.StatusBadRequest)
		return
	}
	learnerID, _ := strconv.ParseInt(fmt.Sprint(r.Context().Value("learner_id")), 10, 32)
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
	messagingSession.Learner.ID = uint(learnerID)
	messagingSession.Pricing = pricing
	//Create
	messagingSessionDAO := daos.GetMessagingSessionDAO()
	result, err := messagingSessionDAO.CreateMessagingSession(messagingSession)
	//reduce learner available coin
	learnerDAO := daos.GetLearnerDAO()
	learner := models.Learner{
		ID:                 uint(learnerID),
		AvailableCoinCount: uint(availableCoinCount) - pricing.Price,
	}
	_, _ = learnerDAO.UpdateLearnerByLearnerID(learner)
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
	if messagingSession.ID != "" {
		http.Error(w, "missing session id.", http.StatusBadRequest)
		return
	} else {
		messagingSession.ID = messagingSessionID
	}
	//Update
	messagingSessionDAO := daos.GetMessagingSessionDAO()
	result, err := messagingSessionDAO.UpdateMessagingSessionByID(messagingSession)

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
	messagingSession.LearnerID = uint(learnerID)
	//parse body
	//Check if user inputs sessionID
	if messagingSession.ID != "" {
		http.Error(w, "missing session id.", http.StatusBadRequest)
		return
	} else {
		messagingSession.ID = messagingSessionID
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
	_, err := messagingSessionDAO.UpdateMessagingSessionByID(messagingSession)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	//Return the coin for learner
	pricingDAO := daos.GetPricingDAO()
	messagingPricing, _ := pricingDAO.GetPricingByID(config.GetPricingIDConfig().MessagingSessionPricingID)
	learnerAvailableCoin, _ := strconv.ParseUint(fmt.Sprint(r.Context().Value("available_coin_count")), 10, 32)
	currentLearner := models.Learner{
		ID:                 uint(learnerID),
		AvailableCoinCount: uint(learnerAvailableCoin) + messagingPricing.Price,
	}
	learnerDAO := daos.GetLearnerDAO()
	result, err := learnerDAO.UpdateLearnerByLearnerID(currentLearner)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, result)

}
