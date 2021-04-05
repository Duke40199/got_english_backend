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

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, result)

}
func UpdateMessagingSession(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		params  = mux.Vars(r)
		message = "OK"
	)
	//parse accountID

	accountID, _ := uuid.Parse(fmt.Sprint(r.Context().Value("id")))
	messagingSession := models.MessagingSession{}
	//parse body
	messagingSessionID := params["messaging_session_id"]
	if err := json.NewDecoder(r.Body).Decode(&messagingSession); err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	//Check if user inputs sessionID
	if messagingSession.ID != "" {
		http.Error(w, "can't update session id.", http.StatusBadRequest)
		return
	}
	messagingSession.ID = messagingSessionID
	learnerDAO := daos.GetLearnerDAO()
	learner, _ := learnerDAO.GetLearnerInfoByAccountID(accountID)

	messagingSessionDAO := daos.GetMessagingSessionDAO()
	messagingSession.LearnerID = learner.ID

	result, err := messagingSessionDAO.UpdateMessagingSessionByID(messagingSession)

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, result)

}
