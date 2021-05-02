package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/models"
	"github.com/golang/got_english_backend/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func CreateTranslationSessionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	//Get learnerID
	accountID, _ := uuid.Parse(fmt.Sprint(r.Context().Value("id")))
	learnerID, _ := strconv.ParseInt(fmt.Sprint(r.Context().Value("learner_id")), 10, 32)
	availableCoinCount, _ := strconv.ParseInt(fmt.Sprint(r.Context().Value("available_coin_count")), 10, 32)
	//Get translation sessions
	translationSession := models.TranslationSession{}
	if err := json.NewDecoder(r.Body).Decode(&translationSession); err != nil {
		http.Error(w, "Malformed data", http.StatusBadRequest)
		return
	}
	if translationSession.ID == "" {
		http.Error(w, "Missing (document) id.", http.StatusBadRequest)
		return
	}
	//Check if have pricing id
	if *translationSession.PricingID == 0 {
		http.Error(w, "Invalid pricing", http.StatusBadRequest)
		return
	}

	//Get pricing
	pricingDAO := daos.GetPricingDAO()
	pricing, _ := pricingDAO.GetPricingByID(*translationSession.PricingID)
	if availableCoinCount < int64(pricing.Price) {
		http.Error(w, "Insufficient coin.", http.StatusBadRequest)
		return
	}
	//add creator learner
	learnerDAO := daos.GetLearnerDAO()
	learner, _ := learnerDAO.GetLearnerInfoByAccountID(accountID)
	//Get exchange rate
	exchangeRateDAO := daos.GetExchangeRateDAO()
	exchangeRate, _ := exchangeRateDAO.GetExchangeRateByServiceName(config.GetServiceConfig().TranslationService)
	translationSession.Learners = append(translationSession.Learners, *learner)
	translationSession.Pricing = pricing
	translationSession.ExchangeRate = *exchangeRate
	translationSession.CreatorLearnerID = uint(learnerID)
	translationSession.CreatedAt = time.Now()
	translationSession.UpdatedAt = time.Now()
	//Create
	translationSessionDAO := daos.GetTranslationSessionDAO()
	result, err := translationSessionDAO.CreateTranslationSession(translationSession)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	//reduce learner available coin
	learner.AvailableCoinCount = uint(availableCoinCount) - pricing.Price
	_, _ = learnerDAO.UpdateLearnerByLearnerID(learner.ID, *learner)
	config.ResponseWithSuccess(w, message, result)

}

func GetTranslationSessionHistoryHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params   = mux.Vars(r)
		message    = "OK"
		timePeriod string
		startDate  time.Time
		endDate    time.Time
		err        error
	)
	learnerID, _ := strconv.ParseUint(fmt.Sprint(r.Context().Value("learner_id")), 10, 0)
	if len(r.URL.Query()["timePeriod"]) > 0 {
		timePeriod = fmt.Sprint(r.URL.Query()["timePeriod"][0])
		startDate, endDate, err = utils.GetTimesByPeriod(timePeriod)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
			return
		}
	} else {
		//Get by month
		startDate, endDate, err = utils.GetTimesByPeriod("monthly")
	}
	translationSessionDAO := daos.GetTranslationSessionDAO()
	result, err := translationSessionDAO.GetTranslationSessionHistory(uint(learnerID), startDate, endDate)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}

func FinishTranslationSessionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		params  = mux.Vars(r)
		message = "OK"
		timeNow = time.Now()
	)
	//parse messagingSessionID
	translationSessionID := params["translation_session_id"]
	//Check if user inputs sessionID
	if translationSessionID == "" {
		http.Error(w, "missing session id.", http.StatusBadRequest)
		return
	}
	//Update session isFinished
	translationSessionDAO := daos.GetTranslationSessionDAO()
	//Check if the session is already finished
	translationSession, _ := translationSessionDAO.GetTranslationSessionByID(translationSessionID)
	if translationSession.IsFinished {
		http.Error(w, "Session is already finished.", http.StatusBadRequest)
		return
	}
	if translationSession.IsCancelled {
		http.Error(w, "Session is already cancelled.", http.StatusBadRequest)
		return
	}
	//Update finish status
	translationSession.IsFinished = true
	translationSession.FinishedAt = &timeNow
	result, translationSession, _ := translationSessionDAO.UpdateTranslationSessionByID(translationSessionID, *translationSession, nil)
	//Get coin value in VND
	pricingDAO := daos.GetPricingDAO()
	coinValue, _ := pricingDAO.GetPricings("coin_value", 0)
	coinValueInVND := (*coinValue)[0].Price
	//Calculate earning
	expertEarnings := utils.CalculateExpertEarningBySession(translationSession.ExchangeRate.Rate, coinValueInVND, translationSession.PaidCoins)
	earningDAO := daos.GetEarningDAO()
	earning := models.Earning{
		Value:                expertEarnings,
		ExpertID:             *translationSession.ExpertID,
		TranslationSessionID: &translationSession.ID,
	}
	_, err := earningDAO.CreateEarning(earning)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}

func UpdateTranslationSessionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		params  = mux.Vars(r)
		message = "OK"
		// learners              = []models.Learner{}
		translationSessionID  = params["translation_session_id"]
		translationSession    = models.TranslationSession{}
		translationSessionDAO = daos.GetTranslationSessionDAO()
	)
	expertID, _ := strconv.ParseUint(fmt.Sprint(r.Context().Value("expert_id")), 10, 0)
	expertIDUint := uint(expertID)
	//if expert is using this endpoint, check if the session already has expert
	if expertID != 0 {
		tmp, err := translationSessionDAO.GetTranslationSessionByID(translationSessionID)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
			return
		}
		if tmp.ExpertID != nil {
			http.Error(w, "An expert is already in this session", http.StatusBadRequest)
			return
		}
		translationSession.ExpertID = &expertIDUint
	} else {
		//Learner is using the endpoint
		//parse body
		if err := json.NewDecoder(r.Body).Decode(&translationSession); err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
			return
		}
		//Check if user inputs sessionID
		if translationSession.ID != "" {
			http.Error(w, "missing session id.", http.StatusBadRequest)
			return
		}
		if translationSession.IsFinished {
			http.Error(w, "Cannot update finish status using this call.", http.StatusBadRequest)
			return
		}
		//Get learners on learnerIDS
		learnerDAO := daos.GetLearnerDAO()
		learners, err := learnerDAO.GetLearnerInfoByIDS(translationSession.LearnerIDs)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
			return
		}
		translationSession.Learners = *learners
	}
	result, _, err := translationSessionDAO.UpdateTranslationSessionByID(translationSessionID, translationSession, translationSession.Learners)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}

func CancelTranslationSessionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		params  = mux.Vars(r)
		message = "OK"
	)
	//parse accountID

	learnerID, _ := strconv.ParseUint(fmt.Sprint(r.Context().Value("learner_id")), 10, 0)
	availableCoinCount, _ := strconv.ParseUint(fmt.Sprint(r.Context().Value("available_coin_count")), 10, 0)
	translationSessionID := params["translation_session_id"]
	if translationSessionID == "" {
		http.Error(w, "missing session id.", http.StatusBadRequest)
		return
	}

	translationSession := models.TranslationSession{
		ID:          translationSessionID,
		IsCancelled: true,
	}
	liveCallSessionDAO := daos.GetTranslationSessionDAO()
	//Check if the session is already cancelled or existed
	tmpSession, _ := liveCallSessionDAO.GetTranslationSessionByID(translationSession.ID)
	if tmpSession.ID == "" {
		http.Error(w, "session not found.", http.StatusBadRequest)
		return
	}
	if tmpSession.IsCancelled {
		http.Error(w, "session is already cancelled.", http.StatusBadRequest)
		return
	}
	result, _, err := liveCallSessionDAO.UpdateTranslationSessionByID(translationSessionID, translationSession, nil)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	//refund
	pricingDAO := daos.GetPricingDAO()
	pricing, _ := pricingDAO.GetPricingByID(*tmpSession.PricingID)
	currentLearner := models.Learner{
		ID:                 uint(learnerID),
		AvailableCoinCount: uint(availableCoinCount) + pricing.Price,
	}
	learnerDAO := daos.GetLearnerDAO()
	_, _ = learnerDAO.UpdateLearnerByLearnerID(uint(learnerID), currentLearner)
	config.ResponseWithSuccess(w, message, result)

}
