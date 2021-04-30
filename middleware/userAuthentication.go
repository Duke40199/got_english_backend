package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
	"github.com/golang/got_english_backend/models"
	"github.com/google/uuid"
)

var roleNameConfig = config.GetRoleNameConfig()

// FirebaseAuthentication : to verify all authorized operations
func UserAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			//Check if token is valid.
			isValidToken, token := CheckIfValidToken(r)
			if !isValidToken {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			claims, _ := token.Claims.(jwt.MapClaims)
			userInfo := claims["claims"].(map[string]interface{})
			//Check if suspended
			_, err := GetAccountFullInfo(userInfo["id"])
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusForbidden)
				return
			}
			//Check if correct role
			if userInfo["role_name"] == "" {
				http.Error(w, "Your current role cannot access this function.", http.StatusForbidden)
				return
			}
			ctx := context.WithValue(r.Context(), "UserAccessToken", token)
			ctx = context.WithValue(ctx, "id", userInfo["id"])
			ctx = context.WithValue(ctx, "role_name", userInfo["role_name"])
			next.ServeHTTP(w, r.WithContext(ctx))

		} else {
			http.Error(w, "Unauthorized", http.StatusForbidden)
		}
	}
}

func ModeratorAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			//Check if token is valid.
			isValidToken, token := CheckIfValidToken(r)
			if !isValidToken {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			claims, _ := token.Claims.(jwt.MapClaims)
			userInfo := claims["claims"].(map[string]interface{})
			//Check if suspended
			accountInfo, err := GetAccountFullInfo(userInfo["id"])
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusForbidden)
				return
			}
			//Check if correct role
			if userInfo["role_name"] != roleNameConfig.Moderator {
				http.Error(w, "Your current role cannot access this function.", http.StatusForbidden)
				return
			}
			ctx := context.WithValue(r.Context(), "UserAccessToken", token)
			//set permission into context
			ctx = context.WithValue(ctx, "can_manage_application_form", accountInfo.Moderator.CanManageApplicationForm)
			ctx = context.WithValue(ctx, "can_manage_coin_bundle", accountInfo.Moderator.CanManageCoinBundle)
			ctx = context.WithValue(ctx, "can_manage_pricing", accountInfo.Moderator.CanManagePricing)
			ctx = context.WithValue(ctx, "can_manage_exchange_rate", accountInfo.Moderator.CanManageExchangeRate)
			ctx = context.WithValue(ctx, "id", userInfo["id"])
			ctx = context.WithValue(ctx, "role_name", userInfo["role_name"])
			next.ServeHTTP(w, r.WithContext(ctx))

		} else {
			http.Error(w, "Unauthorized", http.StatusForbidden)
		}
	}
}
func ModeratorAdminAuthentication(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			//Check if token is valid.
			isValidToken, token := CheckIfValidToken(r)
			if !isValidToken {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			claims, _ := token.Claims.(jwt.MapClaims)
			userInfo := claims["claims"].(map[string]interface{})
			//Check if suspended
			accountInfo, err := GetAccountFullInfo(userInfo["id"])
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusForbidden)
				return
			}
			//Check if correct role
			ctx := context.WithValue(r.Context(), "UserAccessToken", token)
			ctx = context.WithValue(ctx, "id", userInfo["id"])
			ctx = context.WithValue(ctx, "role_name", userInfo["role_name"])
			switch userInfo["role_name"] {
			case roleNameConfig.Moderator:
				{
					//set permission into context
					ctx = context.WithValue(ctx, "can_manage_application_form", accountInfo.Moderator.CanManageCoinBundle)
					ctx = context.WithValue(ctx, "can_manage_coin_bundle", accountInfo.Moderator.CanManageCoinBundle)
					ctx = context.WithValue(ctx, "can_manage_pricing", accountInfo.Moderator.CanManagePricing)

					break
				}
			case roleNameConfig.Admin:
				{
					//set permission into context
					ctx = context.WithValue(ctx, "can_manage_admin", accountInfo.Admin.CanManageAdmin)
					ctx = context.WithValue(ctx, "can_manage_expert", accountInfo.Admin.CanManageExpert)
					ctx = context.WithValue(ctx, "can_manage_moderator", accountInfo.Admin.CanManageModerator)
					ctx = context.WithValue(ctx, "can_manage_learner", accountInfo.Admin.CanManageLearner)
					break
				}
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Unauthorized", http.StatusForbidden)
		}
	}
}
func LearnerAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			//Check if token is valid.
			isValidToken, token := CheckIfValidToken(r)
			if !isValidToken {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			claims, _ := token.Claims.(jwt.MapClaims)
			userInfo := claims["claims"].(map[string]interface{})
			//Check if suspended
			accountInfo, err := GetAccountFullInfo(userInfo["id"])
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusForbidden)
				return
			}
			//Check if correct role
			if userInfo["role_name"] != roleNameConfig.Learner {
				http.Error(w, "Your current role cannot access this function.", http.StatusForbidden)
				return
			}
			//Get permissions and put it in context
			ctx := context.WithValue(r.Context(), "UserAccessToken", token)
			ctx = context.WithValue(ctx, "id", userInfo["id"])
			ctx = context.WithValue(ctx, "learner_id", accountInfo.Learner.ID)
			ctx = context.WithValue(ctx, "available_coin_count", accountInfo.Learner.AvailableCoinCount)
			ctx = context.WithValue(ctx, "role_name", userInfo["role_name"])
			next.ServeHTTP(w, r.WithContext(ctx))

		} else {
			http.Error(w, "Unauthorized", http.StatusForbidden)
		}
	}
}
func LearnerExpertAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			//Check if token is valid.
			isValidToken, token := CheckIfValidToken(r)
			if !isValidToken {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			claims, _ := token.Claims.(jwt.MapClaims)
			userInfo := claims["claims"].(map[string]interface{})
			//Check if suspended
			accountInfo, err := GetAccountFullInfo(userInfo["id"])
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusForbidden)
				return
			}
			//Check if correct role
			if userInfo["role_name"] != roleNameConfig.Learner && userInfo["role_name"] != roleNameConfig.Expert {
				http.Error(w, "Your current role cannot access this function.", http.StatusForbidden)
				return
			}
			ctx := context.WithValue(r.Context(), "UserAccessToken", token)
			ctx = context.WithValue(ctx, "id", userInfo["id"])
			ctx = context.WithValue(ctx, "role_name", userInfo["role_name"])
			switch userInfo["role_name"] {
			case roleNameConfig.Expert:
				{
					//set permission into context
					ctx = context.WithValue(ctx, "expert_id", accountInfo.Expert.ID)
					ctx = context.WithValue(ctx, "can_chat", accountInfo.Expert.CanChat)
					ctx = context.WithValue(ctx, "can_join_live_call_session", accountInfo.Expert.CanJoinLiveCallSession)
					ctx = context.WithValue(ctx, "can_join_translation_session", accountInfo.Expert.CanJoinTranslationSession)
					break
				}
			case roleNameConfig.Learner:
				{
					//set permission into context
					ctx = context.WithValue(ctx, "learner_id", accountInfo.Learner.ID)
					ctx = context.WithValue(ctx, "available_coin_count", accountInfo.Learner.AvailableCoinCount)
					break
				}
			}
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		} else {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
	}
}

// FirebaseAuthentication : to verify all authorized operations
func AdminAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			//Check if token is valid.
			isValidToken, token := CheckIfValidToken(r)
			if !isValidToken {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			claims, _ := token.Claims.(jwt.MapClaims)
			userInfo := claims["claims"].(map[string]interface{})
			//Check if suspended
			accountInfo, err := GetAccountFullInfo(userInfo["id"])
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusForbidden)
				return
			}
			//Check if correct role;
			if userInfo["role_name"] != roleNameConfig.Admin {
				http.Error(w, "Your current role cannot access this function.", http.StatusForbidden)
				return
			}
			ctx := context.WithValue(r.Context(), "UserAccessToken", token)
			ctx = context.WithValue(ctx, "id", userInfo["id"])
			ctx = context.WithValue(ctx, "role_name", userInfo["role_name"])
			//set permission into context
			ctx = context.WithValue(ctx, "admin_id", accountInfo.Admin.ID)
			ctx = context.WithValue(ctx, "can_manage_admin", accountInfo.Admin.CanManageAdmin)
			ctx = context.WithValue(ctx, "can_manage_expert", accountInfo.Admin.CanManageExpert)
			ctx = context.WithValue(ctx, "can_manage_moderator", accountInfo.Admin.CanManageModerator)
			ctx = context.WithValue(ctx, "can_manage_learner", accountInfo.Admin.CanManageLearner)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Unauthorized", http.StatusForbidden)
		}
	}
}

// FirebaseAuthentication : to verify all authorized operations
func ExpertAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			//Check if token is valid.
			isValidToken, token := CheckIfValidToken(r)
			if !isValidToken {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			claims, _ := token.Claims.(jwt.MapClaims)
			userInfo := claims["claims"].(map[string]interface{})
			//Check if suspended
			accountInfo, err := GetAccountFullInfo(userInfo["id"])
			if err != nil {
				http.Error(w, fmt.Sprint(err), http.StatusForbidden)
				return
			}
			//Check if correct role
			if userInfo["role_name"] != roleNameConfig.Expert {
				http.Error(w, "Your current role cannot access this function.", http.StatusForbidden)
				return
			}
			ctx := context.WithValue(r.Context(), "UserAccessToken", token)
			ctx = context.WithValue(ctx, "id", userInfo["id"])
			ctx = context.WithValue(ctx, "role_name", userInfo["role_name"])
			//set permission into context
			ctx = context.WithValue(ctx, "expert_id", accountInfo.Expert.ID)
			ctx = context.WithValue(ctx, "can_chat", accountInfo.Expert.CanChat)
			ctx = context.WithValue(ctx, "can_join_live_call_session", accountInfo.Expert.CanJoinLiveCallSession)
			ctx = context.WithValue(ctx, "can_join_translation_session", accountInfo.Expert.CanJoinTranslationSession)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
	}
}

func CheckIfValidToken(r *http.Request) (bool, *jwt.Token) {
	if r.Header.Get("Authorization") != "" {
		authorizationToken := r.Header.Get("Authorization")
		customToken := strings.TrimSpace(strings.Replace(authorizationToken, "Bearer", "", 1))
		token, _ := jwt.Parse(customToken, nil)
		if token == nil {
			return false, nil
		} else {
			return true, token
		}
	} else {
		return false, nil
	}
}

func GetAccountFullInfo(id interface{}) (*models.Account, error) {
	accountDAO := daos.GetAccountDAO()
	accountID, _ := uuid.Parse(fmt.Sprint(id))
	accountInfo, err := accountDAO.FindAccountByID(accountID)
	if err != nil {
		return nil, err
	}
	//If account doesn't exist.
	if accountInfo == nil {
		return nil, errors.New("account doesn't exist")
	}
	//If account is suspended.
	if accountInfo.IsSuspended {
		return nil, errors.New("your account has been suspended")
	}
	return accountInfo, nil
}
