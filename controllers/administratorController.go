package controllers

import (
	"fmt"
	"net/http"

	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/utils"
)

func GetAdministratorSummaryHandler(w http.ResponseWriter, r *http.Request) {
	var (
		// params   = mux.Vars(r)
		message = "OK"
	)
	result := map[string]interface{}{
		"expert_count":              uint(0),
		"learner_count":             uint(0),
		"messaging_session_count":   uint(0),
		"live_call_session_count":   uint(0),
		"translation_session_count": uint(0),
		"invoice_count":             uint(0),
	}
	//Get time period
	var period string
	if len(r.URL.Query()["period"]) > 0 {
		period = r.URL.Query()["period"][0]
	} else {
		period = "daily"
	}
	startDate, endDate, err := utils.GetTimesByPeriod(period)

	//Get expert count created during the period.
	expertDAO := daos.GetExpertDAO()
	expertCount, err := expertDAO.GetCreatedExpertsInTimePeriod(startDate, endDate)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	result["expert_count"] = expertCount

	//Get created learner count during the period.
	learnerDAO := daos.GetLearnerDAO()
	learnerCount, err := learnerDAO.GetCreatedLearnersInTimePeriod(startDate, endDate)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	result["learner_count"] = learnerCount

	//Get created messaging count during the period.
	messagingSessionDAO := daos.GetMessagingSessionDAO()
	messagingSessionCount, err := messagingSessionDAO.GetCreatedMessagingSessionsInTimePeriod(startDate, endDate)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	result["messaging_session_count"] = messagingSessionCount

	//Get created live call count during the period.
	liveCallSessionDAO := daos.GetLiveCallSessionDAO()
	liveCallSessionCount, err := liveCallSessionDAO.GetCreatedLiveCallSessionsInTimePeriod(startDate, endDate)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	result["live_call_session_count"] = liveCallSessionCount

	//Get created translation call count during the period.
	translationSessionDAO := daos.GetTranslationSessionDAO()
	translationSessionCount, err := translationSessionDAO.GetCreatedTranslationSessionInTimePeriod(startDate, endDate)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	result["translation_session_count"] = translationSessionCount

	//Get created translation call count during the period.
	invoiceDAO := daos.GetInvoiceDAO()
	invoiceCount, err := invoiceDAO.GetCreatedInvoiceInTimePeriod(startDate, endDate)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	result["invoice_count"] = invoiceCount

	config.ResponseWithSuccess(w, message, result)
	//Get count of experts created
}
