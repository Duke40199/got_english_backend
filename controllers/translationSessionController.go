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
	"github.com/google/uuid"
)

func CreateTranslationSessionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	//Get learnerID
	accountID, _ := uuid.Parse(fmt.Sprint(r.Context().Value("id")))
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
	if translationSession.PricingID == 0 {
		http.Error(w, "Invalid pricing", http.StatusBadRequest)
		return
	}

	//Get pricing
	pricingDAO := daos.GetPricingDAO()
	pricing, _ := pricingDAO.GetPricingByID(translationSession.PricingID)
	if availableCoinCount < int64(pricing.Price) {
		http.Error(w, "Insufficient coin.", http.StatusBadRequest)
		return
	}
	learnerDAO := daos.GetLearnerDAO()
	learner, _ := learnerDAO.GetLearnerInfoByAccountID(accountID)
	translationSession.Learners = append(translationSession.Learners, learner)
	translationSession.Pricing = *pricing
	translationSession.CreatedAt = time.Now()
	translationSession.UpdatedAt = time.Now()
	//Create
	translationSessionDAO := daos.GetTranslationSessionDAO()
	result, err := translationSessionDAO.CreateTranslationSession(translationSession)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, result)

}
