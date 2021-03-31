package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	responseConfig "github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func UpdateExpertHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		params  = mux.Vars(r)
		message = "OK"
	)
	accountID, _ := uuid.Parse(params["account_id"])
	expertPermissions := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&expertPermissions); err != nil {
		fmt.Print(err)
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	expertDAO := daos.GetExpertDAO()

	result, err := expertDAO.UpdateExpertByAccountID(accountID, expertPermissions)

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
	}
	responseConfig.ResponseWithSuccess(w, message, result)

}
