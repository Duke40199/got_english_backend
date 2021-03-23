package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang/got_english_backend/config"
	responseConfig "github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/models"
	"github.com/gorilla/mux"
)

func CreateCoinBundleHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	var coinBundle = models.CoinBundle{}

	if err := json.NewDecoder(r.Body).Decode(&coinBundle); err != nil {
		errMsg := "Malformed data"
		responseConfig.ResponseWithError(w, errMsg, err)
	}
	coinBundleDAO := daos.GetCoinBundleDAO()
	result, err := coinBundleDAO.CreateCoinBundle(models.CoinBundle{
		Title:       coinBundle.Title,
		Description: coinBundle.Description,
		Quantity:    uint(coinBundle.Quantity),
		Price:       uint(coinBundle.Price),
		PriceUnit:   coinBundle.PriceUnit,
	},
	)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
	} else {
		responseConfig.ResponseWithSuccess(w, message, result)
	}
}

func GetCoinBundlesHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params   = mux.Vars(r)
		message = "OK"
	)
	coinBundleDAO := daos.GetCoinBundleDAO()
	coinBundles, err := coinBundleDAO.GetCoinBundles()
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
	} else {
		responseConfig.ResponseWithSuccess(w, message, coinBundles)
	}
}

func UpdateCoinBundleHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		message = "OK"
		params  = mux.Vars(r)
	)
	//parse request param to get accountid
	coinBundleID, _ := strconv.ParseInt(params["coin_bundle_id"], 10, 0)
	var coinBundle = models.CoinBundle{
		ID: uint(coinBundleID),
	}
	coinBundleDAO := daos.GetCoinBundleDAO()
	if err := json.NewDecoder(r.Body).Decode(&coinBundle); err != nil {
		errMsg := "Malformed data"
		config.ResponseWithError(w, errMsg, err)
	}
	err := coinBundleDAO.UpdateCoinBundleByID(coinBundle)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
	} else {
		config.ResponseWithSuccess(w, message, 1)
	}
}
