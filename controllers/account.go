package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	responseConfig "github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/models"
	"github.com/golang/got_english_backend/utils"
	"github.com/gorilla/mux"

	"github.com/google/uuid"
)

func CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		message = "OK"
	)
	var account = models.Account{}

	accountDAO := daos.GetAccountDAO()
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		errMsg := "Malformed data"
		responseConfig.ResponseWithError(w, errMsg, err)
	}

	_, err := accountDAO.CreateAccount(models.Account{
		ID:       uuid.New(),
		Username: account.Username,
		Email:    account.Email,
		Password: "123456",
	},
	)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
	} else {
		responseConfig.ResponseWithSuccess(w, message, "Created Successfully.")
	}
}

func UpdateAccountHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		message = "OK"
		params  = mux.Vars(r)
	)
	//parse request param to get accountid
	accountID, _ := uuid.Parse(params["account_id"])
	currentSessionAccountID := r.Context().Value("id")
	//Validate if the account owner is requesting the update.
	if params["account_id"] != currentSessionAccountID {
		http.Error(w, "Only the account owner can update their info.", http.StatusForbidden)
		return
	}
	var account = models.Account{
		ID: accountID,
	}
	accountDAO := daos.GetAccountDAO()
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		errMsg := "Malformed data"
		responseConfig.ResponseWithError(w, errMsg, err)
	}
	err := accountDAO.UpdateAccountByID(account)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
	} else {
		responseConfig.ResponseWithSuccess(w, message, 1)
	}
}

func ViewProfileHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params   = mux.Vars(r)
		message = "OK"
	)
	loginResponse := utils.DecodeFirebaseCustomToken(w, r)
	currentUsername := loginResponse.Username
	accountDAO := daos.GetAccountDAO()
	userDetails, err := accountDAO.FindUserByUsername(models.Account{
		Username: currentUsername,
	})
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
	} else {
		responseConfig.ResponseWithSuccess(w, message, userDetails)
	}
}

func GetAccountsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params   = mux.Vars(r)
		message = "OK"
	)
	var role string
	if len(r.URL.Query()["role"]) > 0 {
		role = r.URL.Query()["role"][0]
	} else {
		role = ""
	}
	var username string
	if len(r.URL.Query()["username"]) > 0 {
		username = r.URL.Query()["username"][0]
	} else {
		username = ""
	}
	accountDAO := daos.GetAccountDAO()
	userDetails, err := accountDAO.GetAccounts(models.Account{
		Username: username,
		RoleName: role,
	})
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
	} else {
		responseConfig.ResponseWithSuccess(w, message, userDetails)
	}
}
