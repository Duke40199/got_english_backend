package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"firebase.google.com/go/auth"
	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
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
	accountID := uuid.New()
	//parsing data
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
	if account.Email == nil || account.Password == nil {
		http.Error(w, "Email or password is invalid or missing", http.StatusBadRequest)
		return
	}
	//Generate username if reqbody doesn't have one
	if account.Username == nil {
		currentTimeMillis := utils.GetCurrentEpochTimeInMiliseconds()
		newUsername := account.RoleName + strconv.FormatInt(currentTimeMillis, 10)
		account.Username = &newUsername
	}
	//create account on db
	result, err := accountDAO.CreateAccount(models.Account{
		ID:       accountID,
		Username: account.Username,
		Email:    account.Email,
		RoleName: account.RoleName,
	}, accountPermission)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	//create account on firebase
	claims := map[string]interface{}{
		"id":        result.ID,
		"email":     result.Email,
		"role_name": result.RoleName,
		"username":  result.Username,
	}
	ctx := r.Context()
	params := (&auth.UserToCreate{}).
		UID(result.ID.String()).
		Email(*result.Email).
		EmailVerified(true).
		DisplayName(*result.Username).
		Password(*account.Password).
		Disabled(false)
	firebaseAuth, context := config.SetupFirebase()
	_, err = firebaseAuth.CreateUser(ctx, params)
	if err != nil {
		fmt.Print("Firebase account already existed.")
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
	updateInfo := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&updateInfo); err != nil {
		http.Error(w, "Malformed data", http.StatusBadRequest)
		return
	}
	//hash password before update
	if updateInfo["password"] != nil {
		hashedPassword, _ := Hash(fmt.Sprint(updateInfo["password"]))
		updateInfo["password"] = hashedPassword
	}
	result, err := accountDAO.UpdateAccountByID(accountID, updateInfo)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	config.ResponseWithSuccess(w, message, result)

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
		Username: &currentUsername,
	})
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
