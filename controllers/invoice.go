package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	responseConfig "github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/models"
)

func CreateInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	invoice := models.Invoice{}
	if err := json.NewDecoder(r.Body).Decode(&invoice); err != nil {
		fmt.Print(err)
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	expertDAO := daos.GetInvoiceDAO()

	result, err := expertDAO.CreateInvoice(invoice)

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
	}
	responseConfig.ResponseWithSuccess(w, message, result)

}
