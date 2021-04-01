package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	responseConfig "github.com/golang/got_english_backend/config"
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
		fmt.Print(err)
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	//Get coin bundle by id
	coinBundleDAO := daos.GetCoinBundleDAO()
	coinBundle, _ := coinBundleDAO.GetCoinBundlesByID(invoice.CoinBundleID)
	if coinBundle.Quantity > learner.AvailableCoinCount {
		errMsg := "Learner does not have sufficient coins"
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
	invoiceDAO := daos.GetInvoiceDAO()
	result, err := invoiceDAO.CreateInvoice(invoice)

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	responseConfig.ResponseWithSuccess(w, message, result)

}
