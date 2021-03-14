package controllers

import (
	"encoding/json"
	"net/http"

	responseConfig "github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/models"
	"github.com/golang/got_english_backend/utils"

	"github.com/google/uuid"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params   = mux.Vars(r)
		message = "OK"
		// username = params["username"]
		// //fullname = params["fullname"]
		// email = params["email"]
	)
	var account = models.Account{}

	userDAO := daos.GetAccountDAO()
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		errMsg := "Malformed data"
		responseConfig.ResponseWithError(w, errMsg, err)
	}

	_, err := userDAO.CreateUser(models.Account{
		ID:       uuid.New(),
		Username: account.Username,
		Email:    account.Email,
		Password: "123456",
	},
	)
	if err != nil {
		panic(err)
	} else {
		resp := map[string]string{"id": "hello"}
		responseConfig.ResponseWithSuccess(w, message, resp)
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
	userDAO := daos.GetAccountDAO()
	userDetails, err := userDAO.FindUserByUsername(models.Account{
		Username: currentUsername,
	})
	if err != nil {
		panic(err)
	} else {
		responseConfig.ResponseWithSuccess(w, message, userDetails)
	}
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
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
	userDAO := daos.GetAccountDAO()
	userDetails, err := userDAO.GetUsers(models.Account{
		Username: username,
		RoleName: role,
	})
	if err != nil {
		panic(err)
	} else {
		responseConfig.ResponseWithSuccess(w, message, userDetails)
	}
}
