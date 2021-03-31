package controllers

import (
	"fmt"
	"net/http"

	responseConfig "github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/models"
	"github.com/google/uuid"
)

func CreateMessagingSessionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	accountID, _ := uuid.Parse(fmt.Sprint(r.Context().Value("id")))
	fmt.Printf("=========================ACCOUNTID:%s", accountID)

	learnerDAO := daos.GetLearnerDAO()
	learner, _ := learnerDAO.GetLearnerInfoByAccountID(accountID)

	messagingSessionDAO := daos.GetMessagingSessionDAO()
	messagingSession := models.MessagingSession{
		LearnerID: learner.ID,
	}

	result, err := messagingSessionDAO.CreateMessagingSession(messagingSession)

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	responseConfig.ResponseWithSuccess(w, message, result)

}
