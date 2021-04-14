package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang/got_english_backend/config"
	responseConfig "github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetExpertsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	var expertID uint64
	var expertDetails interface{}
	expertDAO := daos.GetExpertDAO()
	var err error
	//If the user is looking for another profile
	if len(r.URL.Query()["id"]) > 0 {
		expertID, err = strconv.ParseUint(fmt.Sprint(r.URL.Query()["id"][0]), 10, 0)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
			return
		}
		expertDetails, err = expertDAO.GetExpertByID(uint(expertID))
	} else {
		expertDetails, err = expertDAO.GetExperts()
	}
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, expertDetails)
}

func GetTranslatorExpertsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)

	expertDAO := daos.GetExpertDAO()
	result, err := expertDAO.GetTranslatorExperts()

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	responseConfig.ResponseWithSuccess(w, message, result)

}

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
		return
	}
	responseConfig.ResponseWithSuccess(w, message, result)

}
