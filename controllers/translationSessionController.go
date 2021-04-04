package controllers

import (
	"fmt"
	"net/http"

	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/models"
	"github.com/google/uuid"
)

func CreateTranaltionSessionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	accountID, _ := uuid.Parse(fmt.Sprint(r.Context().Value("id")))

	learnerDAO := daos.GetLearnerDAO()
	learner, _ := learnerDAO.GetLearnerInfoByAccountID(accountID)

	translationSessionDAO := daos.GetTranlsationSessionDAO()
	translationSession := models.TranslationSession{
		Learners:  []*models.Learner{learner},
		PricingID: config.GetPricingIDConfig().MessagingSessionPricingID,
	}

	result, err := translationSessionDAO.CreateTranslationSession(translationSession)

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, result)

}

// func UpdateTranlationSessionHandler(w http.ResponseWriter, r *http.Request) {
// 	defer r.Body.Close()
// 	var (
// 		params  = mux.Vars(r)
// 		message = "OK"
// 	)
// 	//parse accountID

// 	accountID, _ := uuid.Parse(fmt.Sprint(r.Context().Value("id")))
// 	translationSession := models.TranslationSession{}
// 	//parse body
// 	translationSessionID, _ := strconv.ParseInt(params["translation_session_id"], 10, 0)
// 	if err := json.NewDecoder(r.Body).Decode(&translationSession); err != nil {
// 		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
// 		return
// 	}
// 	translationSession.ID = uint(translationSessionID)
// 	learnerDAO := daos.GetLearnerDAO()
// 	learner, _ := learnerDAO.GetLearnerInfoByAccountID(accountID)

// 	translationSessionDAO := daos.GetTranlsationSessionDAO()
// 	translationSession.Learners = learner.ID

// 	result, err := privateCallSessionDAO.UpdatePrivateCallSessionByID(privateCallSession)

// 	if err != nil {
// 		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
// 		return
// 	}
// 	config.ResponseWithSuccess(w, message, result)

// }
