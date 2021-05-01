package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/middleware"
	"github.com/golang/got_english_backend/models"
	"github.com/golang/got_english_backend/utils"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

func CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		message = "OK"
	)

	accountDAO := daos.GetAccountDAO()

	//validations
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Malformed data", http.StatusBadRequest)
		return
	}
	var account = models.Account{}
	if err := json.Unmarshal(requestBody, &account); err != nil {
		http.Error(w, "Malformed account data", http.StatusBadRequest)
		return
	}
	var accountPermission = models.PermissionStruct{}
	if err := json.Unmarshal(requestBody, &accountPermission); err != nil {
		http.Error(w, "Malformed permission data", http.StatusBadRequest)
		return
	}
	if !utils.IsEmailValid(*account.Email) {
		http.Error(w, "Invalid email.", http.StatusBadRequest)
		return
	}
	if account.Password == nil {
		http.Error(w, "Password is invalid or missing", http.StatusBadRequest)
		return
	}
	if account.RoleName == "" {
		http.Error(w, "Missing role_name", http.StatusBadRequest)
		return
	}

	//create account on db
	result, err := accountDAO.CreateAccount(models.Account{
		ID:       account.ID,
		Username: account.Username,
		Email:    account.Email,
		RoleName: account.RoleName,
	}, accountPermission)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	//create account on firebase
	ctx := r.Context()
	params := (&auth.UserToCreate{}).
		UID(result.ID.String()).
		Email(*result.Email).
		EmailVerified(true).
		Password(*account.Password).
		Disabled(false)
	firebaseAuth, context := config.SetupFirebase()
	_, err = firebaseAuth.CreateUser(ctx, params)
	if err != nil {
		fmt.Print("Firebase account already existed.")
	}
	//return firebase custom token with claims
	claims := map[string]interface{}{
		"id":        result.ID,
		"email":     result.Email,
		"role_name": result.RoleName,
		"username":  result.Username,
	}
	token, err := firebaseAuth.CustomTokenWithClaims(context, result.ID.String(), claims)
	if err != nil {
		log.Fatalf("error minting custom token: %v\n", err)
		config.ResponseWithError(w, message, err)
	}
	resp := map[string]interface{}{
		"username": &result.Username,
		"token":    token,
	}
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, resp)
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
	currentSessionRoleName := r.Context().Value("role_name")
	//Validate if the account owner is requesting the update.
	if params["account_id"] != currentSessionAccountID && currentSessionRoleName != config.GetRoleNameConfig().Admin {
		http.Error(w, "Only the account owner or an admin can update this info.", http.StatusForbidden)
		return
	}
	accountDAO := daos.GetAccountDAO()
	//if admin is updating, check admin permission
	accountInfo, err := accountDAO.GetAccountByAccountID(accountID)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	if currentSessionRoleName == config.GetRoleNameConfig().Admin {
		permission := middleware.GetAdminPermissionByRoleName(accountInfo.RoleName)
		isAuthenticated := middleware.CheckAdminPermission(permission, r)
		if !isAuthenticated {
			http.Error(w, "You don't have permission to manage "+strings.ToLower(accountInfo.RoleName)+"s.", http.StatusUnauthorized)
			return
		}
	}

	updateInfo := models.Account{}
	if err := json.NewDecoder(r.Body).Decode(&updateInfo); err != nil {
		http.Error(w, "Malformed data", http.StatusBadRequest)
		return
	}
	if updateInfo.Email != nil {
		http.Error(w, "Cannot update email.", http.StatusBadRequest)
		return
	}
	if updateInfo.Fullname != nil {
		valid, err := utils.IsFullnameValid(*updateInfo.Fullname)
		if !valid {
			http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
			return
		}
	}
	if updateInfo.PhoneNumber != nil {
		valid, err := utils.IsPhoneNumberValid(*updateInfo.PhoneNumber)
		if !valid {
			http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
			return
		}
	}
	if updateInfo.Address != nil {
		valid, err := utils.IsAddressValid(*updateInfo.Address)
		if !valid {
			http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
			return
		}
	}
	if updateInfo.Username != nil {
		valid, err := utils.IsUsernameValid(*updateInfo.Username)
		if !valid {
			http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
			return
		}
	}
	if updateInfo.Birthday != nil {
		valid, err := utils.IsBirthdayValid(*updateInfo.Birthday)
		if !valid {
			http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
			return
		}
	}
	//if password is updated, update it on firebase.
	if updateInfo.Password != nil {
		firebaseUpdateUserParams := (&auth.UserToUpdate{}).
			Password(*updateInfo.Password)
		firebaseAuth, ctx := config.SetupFirebase()
		_, err := firebaseAuth.UpdateUser(ctx, accountID.String(), firebaseUpdateUserParams)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		}
	}
	result, err := accountDAO.UpdateAccountByID(accountID, updateInfo)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	if err != nil {
		fmt.Print("Firebase account already existed.")
	}
	config.ResponseWithSuccess(w, message, result)

}

func SuspendAccountHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		message = "OK"
		params  = mux.Vars(r)
	)
	//parse request param to get accountid
	accountID, _ := uuid.Parse(params["account_id"])
	//Validate if the account owner is requesting the update.
	accountDAO := daos.GetAccountDAO()
	accountToSuspend, err := accountDAO.FindAccountByID(accountID)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	if accountToSuspend.IsSuspended {
		http.Error(w, "Account is already suspended.", http.StatusBadRequest)
		return
	}
	//Check current admin permission
	permission := middleware.GetAdminPermissionByRoleName(accountToSuspend.RoleName)
	isAuthenticated := middleware.CheckAdminPermission(permission, r)
	if !isAuthenticated {
		http.Error(w, "You don't have permission to manage "+accountToSuspend.RoleName+"s.", http.StatusUnauthorized)
		return
	}
	accountInfo, _ := accountDAO.GetAccountByAccountID(accountID)

	if accountInfo.RoleName == config.GetRoleNameConfig().Expert || accountInfo.RoleName == config.GetRoleNameConfig().Learner {
		var learnerID, expertID uint
		if accountInfo.Learner == nil {
			learnerID = 0
		} else {
			learnerID = accountInfo.Learner.ID
		}
		if accountInfo.Expert == nil {
			expertID = 0
		} else {
			expertID = accountInfo.Expert.ID
		}
		//check live call
		liveCallDAO := daos.GetLiveCallSessionDAO()
		inProgress, _ := liveCallDAO.GetLiveCallInProgress(learnerID, expertID)
		if inProgress {
			http.Error(w, "Account is currently in a live call session", http.StatusBadRequest)
			return
		}
		messagingSessionDAO := daos.GetMessagingSessionDAO()
		inProgress, _ = messagingSessionDAO.GetMessagingInProgress(learnerID, expertID)
		if inProgress {
			http.Error(w, "Account is currently in a messaging session", http.StatusBadRequest)
			return
		}
		translationSessionDAO := daos.GetTranslationSessionDAO()
		inProgress, _ = translationSessionDAO.GetTranslationSessionInProgress(learnerID, expertID)
		if inProgress {
			http.Error(w, "Account is currently in a translation session", http.StatusBadRequest)
			return
		}
	}
	result, err := accountDAO.SuspendAccountByID(accountID)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	config.ResponseWithSuccess(w, message, result)

}

func UnsuspendAccountHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		message = "OK"
		params  = mux.Vars(r)
	)
	//parse request param to get accountid
	accountID, _ := uuid.Parse(params["account_id"])
	//Validate if the account owner is requesting the update.
	accountDAO := daos.GetAccountDAO()
	accountToSuspend, err := accountDAO.FindAccountByID(accountID)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	if !accountToSuspend.IsSuspended {
		http.Error(w, "Account is not yet suspended.", http.StatusBadRequest)
		return
	}
	//Check current admin permission
	permission := middleware.GetAdminPermissionByRoleName(accountToSuspend.RoleName)
	isAuthenticated := middleware.CheckAdminPermission(permission, r)
	if !isAuthenticated {
		http.Error(w, "You don't have permission to manage "+accountToSuspend.RoleName+"s.", http.StatusUnauthorized)
		return
	}
	result, err := accountDAO.UnsuspendAccountByID(accountID)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	config.ResponseWithSuccess(w, message, result)

}

func ViewProfileHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params  = mux.Vars(r)
		message = "OK"
	)
	var accountID uuid.UUID
	var err error

	//If the user is looking for another profile
	if len(r.URL.Query()["account_id"]) > 0 {
		accountID, err = uuid.Parse(r.URL.Query()["account_id"][0])
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		}
	} else {
		ctx := r.Context()
		accountID, _ = uuid.Parse(fmt.Sprint(ctx.Value("id")))
	}

	accountDAO := daos.GetAccountDAO()
	userDetails, err := accountDAO.FindAccountByID(accountID)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, userDetails)

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
	permission := middleware.GetAdminPermissionByRoleName(role)
	isAuthenticated := middleware.CheckAdminPermission(permission, r)
	if !isAuthenticated {
		http.Error(w, "You don't have permission to manage "+strings.ToLower(role)+"s.", http.StatusUnauthorized)
		return
	}
	accountDAO := daos.GetAccountDAO()
	userDetails, err := accountDAO.GetAccounts(models.Account{
		Username: &username,
		RoleName: role,
	})
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
	} else {
		config.ResponseWithSuccess(w, message, userDetails)
	}
}

//Hash password
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
