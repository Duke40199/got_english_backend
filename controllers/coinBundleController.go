package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/middleware"
	"github.com/golang/got_english_backend/models"
	"github.com/gorilla/mux"
)

func CreateCoinBundleHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	//Check if current moderator has the manage coin bundle permission
	isPermissioned := middleware.CheckModeratorPermission(config.GetModeratorPermissionConfig().CanManageCoinBundle, r)
	if !isPermissioned {
		http.Error(w, "You don't have permission to manage coin bundles", http.StatusUnauthorized)
		return
	}
	var coinBundle = models.CoinBundle{}
	if err := json.NewDecoder(r.Body).Decode(&coinBundle); err != nil {
		http.Error(w, "Malformed data", http.StatusBadRequest)
		return
	}
	isValidInput, err := models.VaildateCoinBundleInput(coinBundle)
	if !isValidInput {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	coinBundleDAO := daos.GetCoinBundleDAO()
	result, err := coinBundleDAO.CreateCoinBundle(models.CoinBundle{
		Title:       coinBundle.Title,
		Description: coinBundle.Description,
		Quantity:    coinBundle.Quantity,
		PriceUnit:   coinBundle.PriceUnit,
	},
	)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}

func GetCoinBundlesHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message           = "OK"
		queryCoinBundleID = r.URL.Query()["id"]
	)
	ctx := r.Context()
	roleName := fmt.Sprint(ctx.Value("role_name"))
	//If moderator queries, check perm
	if roleName == config.GetRoleNameConfig().Moderator {
		isPermissioned := middleware.CheckModeratorPermission(config.GetModeratorPermissionConfig().CanManageCoinBundle, r)
		if !isPermissioned {
			http.Error(w, "You don't have permission to manage coin bundles", http.StatusUnauthorized)
			return
		}
	}
	//If user input id
	var coinBundleID int64
	var err error
	if len(queryCoinBundleID) > 0 {
		coinBundleID, err = strconv.ParseInt(fmt.Sprint(queryCoinBundleID[0]), 10, 0)
		if err != nil {
			http.Error(w, "Incorrect coin bundle input", http.StatusBadRequest)
			return
		}
	} else {
		coinBundleID = 0
	}
	coinBundleDAO := daos.GetCoinBundleDAO()
	coinBundles, err := coinBundleDAO.GetCoinBundles(uint(coinBundleID))

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, coinBundles)

}

func UpdateCoinBundleHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		message = "OK"
		params  = mux.Vars(r)
	)
	//Check if current moderator has the manage coin bundle permission
	isPermissioned := middleware.CheckModeratorPermission(config.GetModeratorPermissionConfig().CanManageCoinBundle, r)
	if !isPermissioned {
		http.Error(w, "You don't have permission to manage coin bundles", http.StatusUnauthorized)
		return
	}
	//parse request param to get accountid
	coinBundleID, _ := strconv.ParseInt(params["coin_bundle_id"], 10, 0)
	var coinBundle = models.CoinBundle{}
	coinBundleDAO := daos.GetCoinBundleDAO()
	if err := json.NewDecoder(r.Body).Decode(&coinBundle); err != nil {
		http.Error(w, "Malformed data", http.StatusBadRequest)
		return
	}
	isValidInput, err := models.VaildateCoinBundleInput(coinBundle)
	if !isValidInput {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	result, err := coinBundleDAO.UpdateCoinBundleByID(uint(coinBundleID), coinBundle)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}
func DeleteCoinBundleHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		message = "OK"
		params  = mux.Vars(r)
	)
	//Check if current moderator has the manage coin bundle permission
	isPermissioned := middleware.CheckModeratorPermission(config.GetModeratorPermissionConfig().CanManageCoinBundle, r)
	if !isPermissioned {
		http.Error(w, "You don't have permission to manage coin bundles", http.StatusUnauthorized)
		return
	}
	//parse request param to get accountid
	coinBundleID, err := strconv.ParseUint(params["coin_bundle_id"], 10, 0)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	coinBundleDAO := daos.GetCoinBundleDAO()
	result, err := coinBundleDAO.DeleteCoinBundleByID(uint(coinBundleID))
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}
