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

func GetPricingsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	ctx := r.Context()
	roleName := fmt.Sprint(ctx.Value("role_name"))
	//If moderator queries, check perm
	if roleName == config.GetRoleNameConfig().Moderator {
		isPermissioned := middleware.CheckModeratorPermission(config.GetModeratorPermissionConfig().CanManagePricing, r)
		if !isPermissioned {
			http.Error(w, "You don't have permission to manage pricings", http.StatusUnauthorized)
			return
		}
	}
	var pricingID uint = 0
	var pricingName string = ""
	if len(r.URL.Query()["pricing_name"]) > 0 {
		pricingName = fmt.Sprint(r.URL.Query()["pricing_name"][0])
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
	result, err := pricingDAO.GetPricings(pricingName, uint(pricingID))
	if err != nil {
		fmt.Print(err)
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}

func CreatePricingHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		message = "OK"
		// params  = mux.Vars(r)
	)
	//Check if the user is allows to update pricing
	isPermissioned := middleware.CheckModeratorPermission(config.GetModeratorPermissionConfig().CanManagePricing, r)
	if !isPermissioned {
		http.Error(w, "You don't have permission to manage pricings", http.StatusUnauthorized)
		return
	}
	//parse body
	pricing := models.Pricing{}
	if err := json.NewDecoder(r.Body).Decode(&pricing); err != nil {
		http.Error(w, "Malformed data", http.StatusBadRequest)
		return
	}

	pricingDAO := daos.GetPricingDAO()
	result, err := pricingDAO.CreatePricingHandler(pricing)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	config.ResponseWithSuccess(w, message, result)

}

func UpdatePricingHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		message = "OK"
		params  = mux.Vars(r)
	)
	//Check if the user is allows to update pricing
	isPermissioned := middleware.CheckModeratorPermission(config.GetModeratorPermissionConfig().CanManagePricing, r)
	if !isPermissioned {
		http.Error(w, "You don't have permission to manage pricings", http.StatusUnauthorized)
		return
	}
	//parse request param to get accountid
	pricingID, err := strconv.ParseUint(params["pricing_id"], 10, 0)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	//parse body
	updateInfo := models.Pricing{}
	if err := json.NewDecoder(r.Body).Decode(&updateInfo); err != nil {
		http.Error(w, "Malformed data", http.StatusBadRequest)
		return
	}
	pricingDAO := daos.GetPricingDAO()
	result, err := pricingDAO.UpdatePricingByID(uint(pricingID), updateInfo)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	config.ResponseWithSuccess(w, message, result)

}
func DeletePricingHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		message = "OK"
		params  = mux.Vars(r)
	)
	//Check if the user is allows to update pricing
	isPermissioned := middleware.CheckModeratorPermission(config.GetModeratorPermissionConfig().CanManagePricing, r)
	if !isPermissioned {
		http.Error(w, "You don't have permission to manage pricings", http.StatusUnauthorized)
		return
	}
	//parse request param to get accountid
	pricingID, err := strconv.ParseUint(params["pricing_id"], 10, 0)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	pricingDAO := daos.GetPricingDAO()
	result, err := pricingDAO.DeletePricingByID(uint(pricingID))
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}
