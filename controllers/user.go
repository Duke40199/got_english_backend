package controllers

import (
	"encoding/json"
	"net/http"

	responseConfig "github.com/golang/GotEnglishBackend/Application/config"
	"github.com/golang/GotEnglishBackend/Application/daos"
	"github.com/golang/GotEnglishBackend/Application/database"
	"github.com/golang/GotEnglishBackend/Application/models"
	"github.com/golang/GotEnglishBackend/Application/utils"

	"github.com/google/uuid"
)

//CreateUserHandler will create an user and add to DB
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
	db, err := database.ConnectToDB()
	if err == nil {
		userDAO := daos.GetAccountDAO()
		if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
			errMsg := "Malformed data"
			responseConfig.ResponseWithError(w, errMsg, err)
		}

		_, err := userDAO.CreateUser(db, models.Account{
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

}

//ViewProfileHandler will create an user and add to DB
func ViewProfileHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var (
		// params   = mux.Vars(r)
		message = "OK"
	)
	loginResponse := utils.DecodeFirebaseIDToken(w, r)
	currentUsername := loginResponse.Username
	db, err := database.ConnectToDB()
	if err == nil {
		userDAO := daos.GetAccountDAO()
		userDetails, err := userDAO.FindUserByUsername(db, models.Account{
			Username: currentUsername,
		})
		if err != nil {
			panic(err)
		} else {
			responseConfig.ResponseWithSuccess(w, message, userDetails)
		}
	}

}
