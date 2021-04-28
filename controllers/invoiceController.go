package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/models"
	"github.com/google/uuid"
)

func CreateInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	var err error
	//Get current learner
	accountID, _ := uuid.Parse(fmt.Sprint(r.Context().Value("id")))
	learnerDAO := daos.GetLearnerDAO()
	learner, _ := learnerDAO.GetLearnerInfoByAccountID(accountID)
	//Get invoice info
	invoice := models.Invoice{
		ID:      uuid.New(),
		Learner: *learner,
	}
	if err := json.NewDecoder(r.Body).Decode(&invoice); err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	//Get coin bundle by id
	coinBundleDAO := daos.GetCoinBundleDAO()
	coinBundle, err := coinBundleDAO.GetCoinBundleByID(invoice.CoinBundleID, models.CoinBundle{IsDeleted: true})
	if err != nil || coinBundle.ID == 0 {
		http.Error(w, "Coin bundle not found.", http.StatusInternalServerError)
		return
	}
	//Create invoice
	invoiceDAO := daos.GetInvoiceDAO()
	result, err := invoiceDAO.CreateInvoice(invoice)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	learner.AvailableCoinCount += *coinBundle.Quantity
	//update learner available coin after creating invoice.
	_, err = learnerDAO.UpdateLearnerByLearnerID(learner.ID, *learner)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}
