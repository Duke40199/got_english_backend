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

func ValidateManageCoinBundlePermission(w http.ResponseWriter, r *http.Request) bool {
	canManageCoinBundle, _ := strconv.ParseBool(fmt.Sprint(r.Context().Value("can_manage_coin_bundle")))
	if !canManageCoinBundle {
		errMsg := "Your account's permission to 'manage coin bundle' has been disabled."
		http.Error(w, errMsg, http.StatusForbidden)
		return false
	}
	return true
}

func CreateCoinBundleHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	//Check if current moderator has the manage coin bundle permission
	isPermissioned := ValidateManageCoinBundlePermission(w, r)
	if !isPermissioned {
		return
	}
	var coinBundle = models.CoinBundle{}
	if err := json.NewDecoder(r.Body).Decode(&coinBundle); err != nil {
		errMsg := "Malformed data"
		http.Error(w, errMsg, http.StatusBadRequest)
		return
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
		return
	}
	responseConfig.ResponseWithSuccess(w, message, result)
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
		return
	}
	responseConfig.ResponseWithSuccess(w, message, coinBundles)

}

func UpdateCoinBundleHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		message = "OK"
		params  = mux.Vars(r)
	)
	//Check if current moderator has the manage coin bundle permission
	isPermissioned := ValidateManageCoinBundlePermission(w, r)
	if !isPermissioned {
		return
	}
	//parse request param to get accountid
	coinBundleID, _ := strconv.ParseInt(params["coin_bundle_id"], 10, 0)
	var coinBundle = models.CoinBundle{
		ID: uint(coinBundleID),
	}
	coinBundleDAO := daos.GetCoinBundleDAO()
	if err := json.NewDecoder(r.Body).Decode(&coinBundle); err != nil {
		errMsg := "Malformed data"
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	result, err := coinBundleDAO.UpdateCoinBundleByID(coinBundle)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}