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
	accountID, _ := uuid.Parse(fmt.Sprint(r.Context().Value("id")))

	learnerDAO := daos.GetLearnerDAO()
	learner, _ := learnerDAO.GetLearnerInfoByAccountID(accountID)

	messagingSessionDAO := daos.GetMessagingSessionDAO()
	messagingSession := models.MessagingSession{
		LearnerID: learner.ID,
		PricingID: &config.GetPricingIDConfig().MessagingSessionPricingID,
	}

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
	messagingSessionID, _ := strconv.ParseInt(params["messaging_session_id"], 10, 0)
	if err := json.NewDecoder(r.Body).Decode(&messagingSession); err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	messagingSession.ID = uint(messagingSessionID)
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
