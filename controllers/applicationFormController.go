package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/middleware"
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
	//validate if expert has already created a pending form
	applicationFormDAO := daos.GetApplicationFormDAO()
	tmp, _ := applicationFormDAO.GetApplicationForms(models.ApplicationForm{ExpertID: uint(expertID), Status: config.GetApplicationFormStatusConfig().Pending})
	if len(*tmp) > 0 {
		http.Error(w, "you have a pending application form.", http.StatusBadRequest)
		return
	}
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
	result, err := applicationFormDAO.CreateApplicationForm(applicationForm)

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
	}
	config.ResponseWithSuccess(w, message, result)
}

func GetApplicationFormsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params   = mux.Vars(r)
		message = "OK"
		status  string
	)
	ctx := r.Context()
	roleName := fmt.Sprint(ctx.Value("role_name"))
	//If moderator queries, check perm
	if roleName == config.GetRoleNameConfig().Moderator {
		isPermissioned := middleware.CheckModeratorPermission(config.GetModeratorPermissionConfig().CanManageApplicationForm, r)
		if !isPermissioned {
			http.Error(w, "You don't have permission to manage application forms", http.StatusUnauthorized)
			return
		}
	}
	if len(r.URL.Query()["status"]) > 0 {
		status = fmt.Sprint(r.URL.Query()["status"][0])
	}
	applicationForm := models.ApplicationForm{Status: status}
	applicationFormDAO := daos.GetApplicationFormDAO()
	applicationForms, err := applicationFormDAO.GetApplicationForms(applicationForm)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, applicationForms)

}
func GetApplicationFormHistoryHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params   = mux.Vars(r)
		message = "OK"
	)
	expertID, _ := strconv.ParseUint(fmt.Sprint(r.Context().Value("expert_id")), 10, 0)
	applicationFormDAO := daos.GetApplicationFormDAO()
	applicationForms, err := applicationFormDAO.GetApplicationFormsByExpertID(uint(expertID))
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, applicationForms)
}

func ApproveApplicationFormHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		params  = mux.Vars(r)
		message = "OK"
	)
	isPermissioned := middleware.CheckModeratorPermission(config.GetModeratorPermissionConfig().CanManageApplicationForm, r)
	if !isPermissioned {
		http.Error(w, "You don't have permission to manage application forms", http.StatusUnauthorized)
		return
	}
	applicationFormID, err := strconv.ParseUint(fmt.Sprint(params["application_form_id"]), 10, 0)
	if err != nil {
		http.Error(w, "Invalid application form ID (numbers only).", http.StatusBadRequest)
		return
	}
	//check application status
	applicationFormDAO := daos.GetApplicationFormDAO()
	applicationForm, _ := applicationFormDAO.GetApplicationFormByID(uint(applicationFormID))
	if applicationForm.Status != config.GetApplicationFormStatusConfig().Pending {
		http.Error(w, "Application form is already being either approved or rejected.", http.StatusBadRequest)
		return
	}
	//update
	_, err = applicationFormDAO.UpdateApplicationFormByID(uint(applicationFormID), models.ApplicationForm{Status: "Approved"})
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	//update expert permission
	permissions := strings.Split(applicationForm.Type, ",")
	expertDetails := models.Expert{}
	for i := 0; i < len(permissions); i++ {
		switch permissions[i] {
		case "can_chat":
			{
				expertDetails.CanChat = true
				break
			}
		case "can_join_translation_session":
			{
				expertDetails.CanJoinTranslationSession = true
				break
			}
		case "can_join_live_call_session":
			{
				expertDetails.CanJoinLiveCallSession = true
				break
			}
		}
	}
	expertDAO := daos.GetExpertDAO()
	result, err := expertDAO.UpdateExpertByExpertID(applicationForm.ExpertID, expertDetails)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}
func RejectApplicationFormHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		params  = mux.Vars(r)
		message = "OK"
	)
	isPermissioned := middleware.CheckModeratorPermission(config.GetModeratorPermissionConfig().CanManageApplicationForm, r)
	if !isPermissioned {
		http.Error(w, "You don't have permission to manage application forms", http.StatusUnauthorized)
		return
	}
	applicationFormID, err := strconv.ParseUint(fmt.Sprint(params["application_form_id"]), 10, 0)
	if err != nil {
		http.Error(w, "Invalid application form ID (numbers only).", http.StatusBadRequest)
		return
	}
	//check application status
	applicationFormDAO := daos.GetApplicationFormDAO()
	applicationForm, _ := applicationFormDAO.GetApplicationFormByID(uint(applicationFormID))
	if applicationForm.Status != config.GetApplicationFormStatusConfig().Pending {
		http.Error(w, "Application form is already being either approved or rejected.", http.StatusBadRequest)
		return
	}
	//update
	updateResult, err := applicationFormDAO.UpdateApplicationFormByID(uint(applicationFormID), models.ApplicationForm{Status: "Rejected"})
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, updateResult)
}

func DeleteApplicationFormHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		message = "OK"
		params  = mux.Vars(r)
	)
	//parse request param to get accountid
	applicationFormID, err := strconv.ParseUint(params["application_form_id"], 10, 0)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	applicationFormDAO := daos.GetApplicationFormDAO()
	applicationForm, err := applicationFormDAO.GetApplicationFormByID(uint(applicationFormID))
	//Check if expert id matches application form's expert id
	expertID, _ := strconv.ParseUint(fmt.Sprint(r.Context().Value("expert_id")), 10, 0)
	if uint(expertID) != applicationForm.ExpertID {
		http.Error(w, "The queried application form is not existed or is not yours.", http.StatusUnauthorized)
		return
	}
	//Delete application form
	result, err := applicationFormDAO.DeleteApplicationFormByID(uint(applicationFormID))
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	config.ResponseWithSuccess(w, message, result)
}
