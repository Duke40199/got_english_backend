package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	responseConfig "github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func UpdateAdminHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		params  = mux.Vars(r)
		message = "OK"
	)
	canMangeAdmin, _ := strconv.ParseBool(fmt.Sprint(r.Context().Value("can_manage_admin")))
	if !canMangeAdmin {
		http.Error(w, "You don't have the permission to 'manage admin'.", http.StatusForbidden)
		return
	}
	accountID, _ := uuid.Parse(params["account_id"])
	adminPermissions := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&adminPermissions); err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}

	adminDAO := daos.GetAdminDAO()
	result, err := adminDAO.UpdateAdminByAccountID(accountID, adminPermissions)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	responseConfig.ResponseWithSuccess(w, message, result)

}
