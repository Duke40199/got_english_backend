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
	accountID, _ := uuid.Parse(fmt.Sprint(r.Context().Value("id")))

	learnerDAO := daos.GetLearnerDAO()
	learner, _ := learnerDAO.GetLearnerInfoByAccountID(accountID)

	privateCallSessionDAO := daos.GetPrivateCallSessionDAO()
	privateCallSession := models.PrivateCallSession{
		LearnerID: learner.ID,
		PricingID: &config.GetPricingIDConfig().MessagingSessionPricingID,
	}

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
	privateCallSessionID, _ := strconv.ParseInt(params["private_call_session_id"], 10, 0)
	if err := json.NewDecoder(r.Body).Decode(&privateCallSession); err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	privateCallSession.ID = uint(privateCallSessionID)
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
