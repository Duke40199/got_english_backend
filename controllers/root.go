package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	config "github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/models"
)

//LoginHandler will handle the login function
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params   = mux.Vars(r)
		message = "OK"
	)
	var account = models.Account{}
	var result = &models.Account{}
	firebaseAuth, context := config.SetupFirebase()

	accountDAO := daos.GetAccountDAO()
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		errMsg := "Malformed data"
		config.ResponseWithError(w, errMsg, err)
	}
	if account.Email == "" {
		result, _ = accountDAO.FindAccountByUsernameAndPassword(account)
	} else {
		result, _ = accountDAO.FindAccountByEmailAndPassword(account)
	}
	//account not found.
	if result.Username == "" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	//set account role for token
	claims := map[string]interface{}{
		"id":        result.ID,
		"role_name": result.RoleName,
		"username":  &result.Username,
	}
	token, err := firebaseAuth.CustomTokenWithClaims(context, "firebase_UID", claims)
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
