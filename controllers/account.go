package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	var accountInfo = daos.AccountFullInfo{}
	accountDAO := daos.GetAccountDAO()
	if err := json.NewDecoder(r.Body).Decode(&accountInfo); err != nil {
		errMsg := "Malformed data"
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	fmt.Print(accountInfo.Email)
	if (accountInfo.Email) == "" {
		http.Error(w, "Email invalid or missing", http.StatusBadRequest)
		return
	}
	accountID := uuid.New()
	_, err := accountDAO.CreateAccount(models.Account{
		ID:       accountID,
		Username: accountInfo.Username,
		Email:    accountInfo.Email,
		Password: accountInfo.Password,
	},
	)
	//add role specific info
	switch accountInfo.RoleName {
	case config.GetRoleNameConfig().Admin:
		{
			admin := models.Admin{
				CanManageExpert:  accountInfo.CanManageExpert,
				CanManageLearner: accountInfo.CanManageLearner,
				CanManageAdmin:   accountInfo.CanManageAdmin,
				AccountID:        accountID,
			}
			adminDAO := daos.GetAdminDAO()
			_, err = adminDAO.CreateAdmin(admin)
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
				return
			}
			break
		}
	case config.GetRoleNameConfig().Expert:
		{
			expert := models.Expert{
				CanChat:                   accountInfo.CanChat,
				CanJoinTranslationSession: accountInfo.CanJoinTranslationSession,
				CanJoinPrivateCallSession: accountInfo.CanJoinPrivateCallSession,
				AccountID:                 accountID,
			}
			expertDAO := daos.GetExpertDAO()
			_, err = expertDAO.CreateExpert(expert)
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
				return
			}
			break
		}
	case config.GetRoleNameConfig().Learner:
		{
			learner := models.Learner{
				AvailableCoinCount: 0,
				AccountID:          accountID,
			}
			learnerDAO := daos.GetLearnerDAO()
			_, err = learnerDAO.CreateLearner(learner)
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
				return
			}
			break
		}
	case config.GetRoleNameConfig().Moderator:
		{
			moderator := models.Moderator{
				CanManageCoinBundle:      accountInfo.CanManageCoinBundle,
				CanManagePricing:         accountInfo.CanManagePricing,
				CanManageApplicationForm: accountInfo.CanManageApplicationForm,
				AccountID:                accountID,
			}
			moderatorDAO := daos.GetModeratorDAO()
			_, err = moderatorDAO.CreateModerator(moderator)
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
				return
			}
			break
		}
	}
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	config.ResponseWithSuccess(w, message, accountID)
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
		errMsg := "Malformed data"
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	//hash password before update
	if updateInfo["password"] != "" {
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
		Username: currentUsername,
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
		Username: username,
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
