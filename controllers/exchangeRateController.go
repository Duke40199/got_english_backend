package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/models"
	"github.com/gorilla/mux"
)

func ValidateManageExchangeRatePermission(w http.ResponseWriter, r *http.Request) bool {
	canManageExchangeRate, _ := strconv.ParseBool(fmt.Sprint(r.Context().Value("can_manage_exchange_rate")))
	if !canManageExchangeRate {
		http.Error(w, "You don't have the permission to manage exchange rates", http.StatusUnauthorized)
		return false
	}
	return true
}

func GetExchangeRatesHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	//Get perm
	if !ValidateManageExchangeRatePermission(w, r) {
		return
	}
	exchangeRate := models.ExchangeRate{}
	//Get id
	if len(r.URL.Query()["id"]) > 0 {
		id, err := strconv.ParseUint(fmt.Sprint(r.URL.Query()["id"][0]), 10, 0)
		if err != nil {
			http.Error(w, "invalid exchange rate id ", http.StatusBadGateway)
			return
		}
		exchangeRate.ID = uint(id)
	}
	if len(r.URL.Query()["service_name"]) > 0 {
		serviceName := fmt.Sprint(r.URL.Query()["service_name"][0])
		exchangeRate.ServiceName = serviceName
	}
	exchangeRateDAO := daos.GetExchangeRateDAO()
	result, err := exchangeRateDAO.GetExchangeRates(exchangeRate)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}

func UpdateExchangeRateHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		params  = mux.Vars(r)
		message = "OK"
	)
	if !ValidateManageExchangeRatePermission(w, r) {
		return
	}
	exchangeRate := models.ExchangeRate{}
	exchangeRateID, err := strconv.ParseInt(params["exchange_rate_id"], 10, 0)
	if err != nil {
		http.Error(w, "Invalid exchange rate id", http.StatusBadRequest)
		return
	}
	exchangeRate.ID = uint(exchangeRateID)
	if err := json.NewDecoder(r.Body).Decode(&exchangeRate); err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	if exchangeRate.Rate < 0 || exchangeRate.Rate > 1 {
		http.Error(w, "Rate is between 0 and 1", http.StatusBadRequest)
		return
	}
	exchangeRateDAO := daos.GetExchangeRateDAO()
	result, err := exchangeRateDAO.UpdateExchangeRateByID(exchangeRate.ID, exchangeRate)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}
