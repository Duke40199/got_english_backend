package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/utils"
)

func GetAdministratorSummaryHandler(w http.ResponseWriter, r *http.Request) {
	var (
		// params   = mux.Vars(r)
		message = "OK"
	)
	result := map[string]interface{}{}
	startDateDaily, endDateDaily, _ := utils.GetTimesByPeriod("daily")
	startDateWeekly, _, _ := utils.GetTimesByPeriod("weekly")
	//Get expert count created during the period.
	expertDAO := daos.GetExpertDAO()
	expertCount, err := expertDAO.GetCreatedExpertsInTimePeriod(startDateDaily, endDateDaily)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	result["new_expert_daily_count"] = expertCount

	//Expert weekly
	expertWeeklyCount, err := expertDAO.GetNewExpertsCountInTimePeriod(startDateWeekly.AddDate(0, 0, -7), startDateWeekly)
	result["new_expert_weekly_count"] = expertWeeklyCount

	//Get created learner count during the period.
	learnerDAO := daos.GetLearnerDAO()
	learnerCount, err := learnerDAO.GetCreatedLearnersInTimePeriod(startDateDaily, endDateDaily)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	result["new_learner_daily_count"] = learnerCount

	//Learner weekly
	learnerWeeklyCount, err := learnerDAO.GetNewLearnersCountInTimePeriod(startDateWeekly.AddDate(0, 0, -7), startDateWeekly)
	result["new_learner_weekly_count"] = learnerWeeklyCount

	//Get created messaging count during the period.
	messagingSessionDAO := daos.GetMessagingSessionDAO()
	messagingSessionCount, err := messagingSessionDAO.GetCreatedMessagingSessionsInTimePeriod(startDateDaily, endDateDaily)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	result["new_messaging_session_daily_count"] = messagingSessionCount
	//messaging session weekly
	messagingSessionWeeklyCount, err := messagingSessionDAO.GetNewMessagingSessionsCountInTimePeriod(startDateWeekly.AddDate(0, 0, -7), startDateWeekly)
	result["new_messaging_session_weekly_count"] = messagingSessionWeeklyCount

	//Get created live call count during the period.
	liveCallSessionDAO := daos.GetLiveCallSessionDAO()
	liveCallSessionCount, err := liveCallSessionDAO.GetCreatedLiveCallSessionsInTimePeriod(startDateDaily, endDateDaily)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	result["new_live_call_session_daily_count"] = liveCallSessionCount

	//messaging session weekly
	liveCallSessionsWeeklyCount, err := liveCallSessionDAO.GetNewLiveCallSessionsCountInTimePeriod(startDateWeekly.AddDate(0, 0, -7), startDateWeekly)
	result["new_live_call_session_weekly_count"] = liveCallSessionsWeeklyCount

	//Get created translation call count during the period.
	translationSessionDAO := daos.GetTranslationSessionDAO()
	translationSessionCount, err := translationSessionDAO.GetCreatedTranslationSessionInTimePeriod(startDateDaily, endDateDaily)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	result["new_translation_session_daily_count"] = translationSessionCount

	//messaging session weekly
	translationSessionsWeeklyCount, err := translationSessionDAO.GetNewTranslationSessionsCountInTimePeriod(startDateWeekly.AddDate(0, 0, -7), startDateWeekly)
	result["new_translation_session_weekly_count"] = translationSessionsWeeklyCount

	//Get created translation call count during the period.
	invoiceDAO := daos.GetInvoiceDAO()
	invoiceCount, err := invoiceDAO.GetCreatedInvoiceInTimePeriod(startDateDaily, endDateDaily)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	result["new_invoice_daily_count"] = invoiceCount
	//messaging session weekly
	invoiceWeeklyCount, err := invoiceDAO.GetNewInvoicesCountInTimePeriod(startDateWeekly.AddDate(0, 0, -7), startDateWeekly)
	result["new_invoice_weekly_count"] = invoiceWeeklyCount

	config.ResponseWithSuccess(w, message, result)
	//Get count of experts created
}

func GetAdministratorMonthlyServiceSummaryHandler(w http.ResponseWriter, r *http.Request) {
	var (
		// params   = mux.Vars(r)
		message = "OK"
		month   time.Month
		year    int
		timeNow = time.Now()
	)
	result := map[string]interface{}{}
	//Get time period
	if len(r.URL.Query()["month"]) > 0 {
		tmp, err := strconv.ParseUint(fmt.Sprint(r.URL.Query()["month"][0]), 10, 0)
		if err != nil || tmp < 1 || tmp > 12 {
			http.Error(w, "Invalid month input", http.StatusBadRequest)
			return
		}
		month = time.Month(tmp)
	} else {
		month = timeNow.Month()
	}
	//Get time period
	if len(r.URL.Query()["year"]) > 0 {
		tmp, err := strconv.ParseUint(fmt.Sprint(r.URL.Query()["year"][0]), 10, 0)
		if err != nil {
			http.Error(w, "Invalid year input", http.StatusBadRequest)
			return
		}
		year = int(tmp)
	} else {
		year = timeNow.Year()
	}
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	lastDay := time.Date(year, month+1, 1, 0, 0, 0, -1, time.Local)
	//messaging session weekly
	messagingSessionDAO := daos.GetMessagingSessionDAO()
	messagingSessionWeeklyCount, err := messagingSessionDAO.GetNewMessagingSessionsCountInTimePeriod(firstDay, lastDay)
	result["new_messaging_session_monthly_count"] = messagingSessionWeeklyCount
	//messaging session weekly
	liveCallSessionDAO := daos.GetLiveCallSessionDAO()
	liveCallSessionsWeeklyCount, err := liveCallSessionDAO.GetNewLiveCallSessionsCountInTimePeriod(firstDay, lastDay)
	result["new_live_call_session_monthly_count"] = liveCallSessionsWeeklyCount
	//messaging session weekly
	translationSessionDAO := daos.GetTranslationSessionDAO()
	translationSessionsWeeklyCount, err := translationSessionDAO.GetNewTranslationSessionsCountInTimePeriod(firstDay, lastDay)
	result["new_translation_session_monthly_count"] = translationSessionsWeeklyCount
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}

func GetAdministratorMonthlyAccountSummaryHandler(w http.ResponseWriter, r *http.Request) {
	var (
		// params   = mux.Vars(r)
		message = "OK"
		month   time.Month
		year    int
		timeNow = time.Now()
	)
	result := map[string]interface{}{}
	//Get time period
	if len(r.URL.Query()["month"]) > 0 {
		tmp, err := strconv.ParseUint(fmt.Sprint(r.URL.Query()["month"][0]), 10, 0)
		if err != nil || tmp < 1 || tmp > 12 {
			http.Error(w, "Invalid month input", http.StatusBadRequest)
			return
		}
		month = time.Month(tmp)
	} else {
		month = timeNow.Month()
	}
	//Get time period
	if len(r.URL.Query()["year"]) > 0 {
		tmp, err := strconv.ParseUint(fmt.Sprint(r.URL.Query()["year"][0]), 10, 0)
		if err != nil {
			http.Error(w, "Invalid year input", http.StatusBadRequest)
			return
		}
		year = int(tmp)
	} else {
		year = timeNow.Year()
	}
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	lastDay := time.Date(year, month+1, 1, 0, 0, 0, -1, time.Local)
	//messaging session weekly
	learnerDAO := daos.GetLearnerDAO()
	newLearnerMonthlyCount, err := learnerDAO.GetNewLearnersCountInTimePeriod(firstDay, lastDay)
	result["new_learner_monthly_count"] = newLearnerMonthlyCount
	//messaging session weekly
	expertDAO := daos.GetExpertDAO()
	newExpertMontlyCount, err := expertDAO.GetNewExpertsCountInTimePeriod(firstDay, lastDay)
	result["new_expert_monthly_count"] = newExpertMontlyCount
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}
