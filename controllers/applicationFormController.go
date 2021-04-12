package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	responseConfig "github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/models"
	"github.com/gorilla/mux"
)

func CreateApplicationFormHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	applicationForm := models.ApplicationForm{}
	expertID, _ := strconv.ParseUint(fmt.Sprint(r.Context().Value("expert_id")), 10, 0)
	if err := json.NewDecoder(r.Body).Decode(&applicationForm); err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	if len(*applicationForm.Types) == 0 {
		http.Error(w, "Missing form types", http.StatusBadRequest)
		return
	}
	for i := 0; i < len(*applicationForm.Types); i++ {
		formType := (*applicationForm.Types)[i]
		if formType != "can_chat" && formType != "can_join_translation_session" && formType != "can_join_live_call_session" {
			http.Error(w, "incorrect application type. (can_chat|can_join_translation_session|can_join_live_call_session)", http.StatusBadRequest)
			return
		}
		applicationForm.Type += formType
		if i < len(*applicationForm.Types)-1 {
			applicationForm.Type += ","
		}
	}
	//assign expertID to applicationForm
	applicationForm.ExpertID = uint(expertID)
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
func ApproveApplicationFormHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		params  = mux.Vars(r)
		message = "OK"
	)
	canManageApplicationForm, _ := strconv.ParseBool(fmt.Sprint(r.Context().Value("can_manage_application_form")))
	if !canManageApplicationForm {
		http.Error(w, "You don't have the permission to manage application forms.", http.StatusUnauthorized)
		return
	}
	applicationFormID, err := strconv.ParseUint(fmt.Sprint(params["application_form_id"]), 10, 0)
	if err != nil {
		http.Error(w, "Invalid application form ID (numbers only).", http.StatusBadRequest)
		return
	}
	//check application status
	applicationFormDAO := daos.GetApplicationFormDAO()
	applicationForm, err := applicationFormDAO.GetApplicationFormByID(uint(applicationFormID))
	if applicationForm.Status != responseConfig.GetApplicationFormStatusConfig().Pending {
		http.Error(w, "Application form is already being either approved or rejected.", http.StatusBadRequest)
		return
	}
	//update
	updateResult, err := applicationFormDAO.UpdateApplicationFormByID(uint(applicationFormID), models.ApplicationForm{Status: "Approved"})
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	responseConfig.ResponseWithSuccess(w, message, updateResult)
}
func RejectApplicationFormHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		params  = mux.Vars(r)
		message = "OK"
	)
	canManageApplicationForm, _ := strconv.ParseBool(fmt.Sprint(r.Context().Value("can_manage_application_form")))
	if !canManageApplicationForm {
		http.Error(w, "You don't have the permission to manage application forms.", http.StatusUnauthorized)
		return
	}
	applicationFormID, err := strconv.ParseUint(fmt.Sprint(params["application_form_id"]), 10, 0)
	if err != nil {
		http.Error(w, "Invalid application form ID (numbers only).", http.StatusBadRequest)
		return
	}
	//check application status
	applicationFormDAO := daos.GetApplicationFormDAO()
	applicationForm, err := applicationFormDAO.GetApplicationFormByID(uint(applicationFormID))
	if applicationForm.Status != responseConfig.GetApplicationFormStatusConfig().Pending {
		http.Error(w, "Application form is already being either approved or rejected.", http.StatusBadRequest)
		return
	}
	//update
	updateResult, err := applicationFormDAO.UpdateApplicationFormByID(uint(applicationFormID), models.ApplicationForm{Status: "Rejected"})
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	responseConfig.ResponseWithSuccess(w, message, updateResult)
}
