package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	responseConfig "github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/models"
	"github.com/google/uuid"
)

func CreateApplicationFormHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	applicationForm := models.ApplicationForm{}
	currentAccountID, _ := uuid.Parse(fmt.Sprint(r.Context().Value("id")))
	if err := json.NewDecoder(r.Body).Decode(&applicationForm); err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	//Get expertID
	expertDAO := daos.GetExpertDAO()
	expert, err := expertDAO.GetExpertByAccountID(currentAccountID)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	//assign expertID to applicationForm
	applicationForm.ExpertID = expert.ID
	applicationFormDAO := daos.GetApplicationFormDAO()
	result, err := applicationFormDAO.CreateApplicationForm(applicationForm)

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
	}
	responseConfig.ResponseWithSuccess(w, message, result)
}

func GetApplicationFormsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params   = mux.Vars(r)
		message = "OK"
	)
	applicationFormDAO := daos.GetApplicationFormDAO()
	applicationForms, err := applicationFormDAO.GetApplicationForms()
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	responseConfig.ResponseWithSuccess(w, message, applicationForms)

}
