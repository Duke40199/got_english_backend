package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/models"
)

func UpdateExchangeRateHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	exchangeRate := models.ExchangeRate{}
	//Get perm
	canManageExchangeRatePermission, _ := strconv.ParseBool(fmt.Sprint(r.Context().Value("can_manage_exchange_rate")))
	if !canManageExchangeRatePermission {
		http.Error(w, "You cannot manage exchange rate", http.StatusUnauthorized)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&exchangeRate); err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	if exchangeRate.ID == 0 {
		http.Error(w, "missing or invalid id", http.StatusBadRequest)
		return
	}
	exchangeRateDAO := daos.GetExchangeRateDAO()
	result, err := exchangeRateDAO.UpdateExchangeRate(exchangeRate)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}
