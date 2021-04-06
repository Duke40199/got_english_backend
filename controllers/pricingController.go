package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
)

func GetPricingsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	var pricingID uint = 0
	var serviceName string = ""
	if len(r.URL.Query()["service_name"]) > 0 {
		serviceName = fmt.Sprint(r.URL.Query()["service_name"][0])
	}
	if len(r.URL.Query()["id"]) > 0 {
		tmp, err := strconv.ParseUint(fmt.Sprint(r.URL.Query()["id"][0]), 10, 0)
		if err != nil {
			fmt.Print(err)
			http.Error(w, "invalid pricing id.", http.StatusBadRequest)
			return
		}
		pricingID = uint(tmp)
	}
	pricingDAO := daos.GetPricingDAO()
	result, err := pricingDAO.GetPricings(serviceName, uint(pricingID))
	if err != nil {
		fmt.Print(err)
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}
