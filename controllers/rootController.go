package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"firebase.google.com/go/auth"
	config "github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/models"
	"github.com/golang/got_english_backend/utils"
	"github.com/google/uuid"
)

//LoginHandler will handle the login function
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params   = mux.Vars(r)
		message = "OK"
	)
	accountDAO := daos.GetAccountDAO()
	var result = &models.Account{}
	var err error
	firebaseAuth, context := config.SetupFirebase()

	//Login with firebase's IDToken
	if r.Header.Get("Authorization") != "" {
		firebaseIDTokenStruct := utils.DecodeFirebaseIDToken(w, r)
		if firebaseIDTokenStruct.Email == "" {
			http.Error(w, "Invalid or expired ID token.", http.StatusForbidden)
			return
		}
		result, err = accountDAO.FindAccountByEmail(models.Account{Email: &firebaseIDTokenStruct.Email})
		if err != nil {
			http.Error(w, fmt.Sprint(err.Error()), http.StatusForbidden)
			return
		}
	} else {
		http.Error(w, "Missing Firebase ID Token.", http.StatusForbidden)
		return
	}
	//account not found.
	if result.Email == nil {
		http.Error(w, "Wrong username/email or password.", http.StatusForbidden)
		return
	}
	//set account role for token
	claims := map[string]interface{}{
		"id":        result.ID,
		"email":     result.Email,
		"role_name": result.RoleName,
		"username":  &result.Username,
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
	config.ResponseWithSuccess(w, message, resp)
}

func LoginWithGoogleHandler(w http.ResponseWriter, r *http.Request) {
	var (
		// params   = mux.Vars(r)
		message = "OK"
	)
	var account = models.Account{}
	var result = &models.Account{}
	firebaseAuth, context := config.SetupFirebase()

	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, "Malformed data", http.StatusBadRequest)
		return
	}
	//Get Google IDToken
	decodedIDToken := utils.DecodeGoogleToken(w, r)
	if decodedIDToken.Email == "" {
		http.Error(w, "Invalid or expired id_token", http.StatusForbidden)
		return
	}
	//Create a user based on google IDToken
	currentTimeMillis := utils.GetCurrentEpochTimeInMiliseconds()
	//Remove spaces
	var generatedUsername = strings.ReplaceAll(decodedIDToken.GivenName+decodedIDToken.FamilyName+strconv.FormatInt(currentTimeMillis, 10), " ", "")
	account = models.Account{
		ID:        uuid.New(),
		Username:  &generatedUsername,
		Email:     &decodedIDToken.Email,
		AvatarURL: &decodedIDToken.Picture,
		Fullname:  &decodedIDToken.Name,
		RoleName:  account.RoleName,
	}
	//If learner is logging in using google for the first time
	//Create a new account from firebase to db
	accountDAO := daos.GetAccountDAO()
	result, _ = accountDAO.FindAccountByEmail(account)
	if result.Email == nil {
		result, _ = accountDAO.CreateAccount(account, models.PermissionStruct{})
	}
	//set account role for token
	claims := map[string]interface{}{
		"id":        result.ID,
		"email":     result.Email,
		"role_name": result.RoleName,
		"username":  result.Username,
	}
	ctx := r.Context()
	//if login for the first time
	params := (&auth.UserToCreate{}).
		UID(result.ID.String()).
		Email(*result.Email).
		EmailVerified(true).
		DisplayName(*result.Username).
		Disabled(false)
	_, err := firebaseAuth.CreateUser(ctx, params)
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
	config.ResponseWithSuccess(w, message, resp)
}
