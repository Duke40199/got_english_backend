package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/daos"
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
			isSuspended := CheckIfSuspended(userInfo["id"])
			if isSuspended {
				http.Error(w, "Your account has been suspended.", http.StatusUnauthorized)
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
			isSuspended := CheckIfSuspended(userInfo["id"])
			if isSuspended {
				http.Error(w, "Your account has been suspended.", http.StatusUnauthorized)
				return
			}
			//Check if correct role
			if userInfo["role_name"] != roleNameConfig.Moderator {
				http.Error(w, "Your current role cannot access this function.", http.StatusForbidden)
				return
			}
			//Get permissions and put it in context
			moderatorDAO := daos.GetModeratorDAO()
			accountID, _ := uuid.Parse(fmt.Sprint(userInfo["id"]))
			permissions, _ := moderatorDAO.GetModeratorByAccountID(accountID)

			ctx := context.WithValue(r.Context(), "UserAccessToken", token)
			//set permission into context
			ctx = context.WithValue(ctx, "can_manage_application_form", permissions.CanManageCoinBundle)
			ctx = context.WithValue(ctx, "can_manage_coin_bundle", permissions.CanManageCoinBundle)
			ctx = context.WithValue(ctx, "can_manage_pricing", permissions.CanManagePricing)
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
			isSuspended := CheckIfSuspended(userInfo["id"])
			if isSuspended {
				http.Error(w, "Your account has been suspended.", http.StatusUnauthorized)
				return
			}
			//Check if correct role
			ctx := context.WithValue(r.Context(), "UserAccessToken", token)
			ctx = context.WithValue(ctx, "id", userInfo["id"])
			ctx = context.WithValue(ctx, "role_name", userInfo["role_name"])
			switch userInfo["role_name"] {
			case roleNameConfig.Moderator:
				{
					//Get permissions and put it in context
					moderatorDAO := daos.GetModeratorDAO()
					accountID, _ := uuid.Parse(fmt.Sprint(userInfo["id"]))
					permissions, _ := moderatorDAO.GetModeratorByAccountID(accountID)
					//set permission into context
					ctx = context.WithValue(ctx, "can_manage_application_form", permissions.CanManageCoinBundle)
					ctx = context.WithValue(ctx, "can_manage_coin_bundle", permissions.CanManageCoinBundle)
					ctx = context.WithValue(ctx, "can_manage_pricing", permissions.CanManagePricing)

					break
				}
			case roleNameConfig.Admin:
				{
					//Get permissions and put it in context
					adminDAO := daos.GetAdminDAO()
					accountID, _ := uuid.Parse(fmt.Sprint(userInfo["id"]))
					permissions, _ := adminDAO.GetAdminByAccountID(accountID)
					//set permission into context
					ctx = context.WithValue(ctx, "can_manage_admin", permissions.CanManageAdmin)
					ctx = context.WithValue(ctx, "can_manage_expert", permissions.CanManageExpert)
					ctx = context.WithValue(ctx, "can_manage_moderator", permissions.CanManageModerator)
					ctx = context.WithValue(ctx, "can_manage_learner", permissions.CanManageLearner)
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
			isSuspended := CheckIfSuspended(userInfo["id"])
			if isSuspended {
				http.Error(w, "Your account has been suspended.", http.StatusUnauthorized)
				return
			}
			//Check if correct role
			if userInfo["role_name"] != roleNameConfig.Learner {
				http.Error(w, "Your current role cannot access this function.", http.StatusForbidden)
				return
			}
			//Get permissions and put it in context
			learnerDAO := daos.GetLearnerDAO()
			accountID, _ := uuid.Parse(fmt.Sprint(userInfo["id"]))
			learnerInfo, _ := learnerDAO.GetLearnerInfoByAccountID(accountID)
			ctx := context.WithValue(r.Context(), "UserAccessToken", token)
			ctx = context.WithValue(ctx, "id", userInfo["id"])
			ctx = context.WithValue(ctx, "learner_id", learnerInfo.ID)
			ctx = context.WithValue(ctx, "available_coin_count", learnerInfo.AvailableCoinCount)
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
			isSuspended := CheckIfSuspended(userInfo["id"])
			if isSuspended {
				http.Error(w, "Your account has been suspended.", http.StatusUnauthorized)
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
					//Get permissions and put it in context
					expertDAO := daos.GetExpertDAO()
					accountID, _ := uuid.Parse(fmt.Sprint(userInfo["id"]))
					permissions, _ := expertDAO.GetExpertByAccountID(accountID)
					//set permission into context
					ctx = context.WithValue(ctx, "can_chat", permissions.CanChat)
					ctx = context.WithValue(ctx, "can_join_live_call_session", permissions.CanJoinLiveCallSession)
					ctx = context.WithValue(ctx, "can_join_translation_room", permissions.CanJoinTranslationSession)
					next.ServeHTTP(w, r.WithContext(ctx))
					break
				}
			case roleNameConfig.Learner:
				{
					//Get permissions and put it in context
					learnerDAO := daos.GetLearnerDAO()
					accountID, _ := uuid.Parse(fmt.Sprint(userInfo["id"]))
					learnerInfo, _ := learnerDAO.GetLearnerInfoByAccountID(accountID)
					//set permission into context
					ctx = context.WithValue(ctx, "learner_id", learnerInfo.ID)
					ctx = context.WithValue(ctx, "available_coin_count", learnerInfo.AvailableCoinCount)
					break
				}
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Unauthorized", http.StatusForbidden)
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
			isSuspended := CheckIfSuspended(userInfo["id"])
			if isSuspended {
				http.Error(w, "Your account has been suspended.", http.StatusUnauthorized)
				return
			}
			//Check if correct role;
			if userInfo["role_name"] != roleNameConfig.Admin {
				http.Error(w, "Your current role cannot access this function.", http.StatusForbidden)
				return
			}
			//Get permissions and put it in context
			adminDAO := daos.GetAdminDAO()
			accountID, _ := uuid.Parse(fmt.Sprint(userInfo["id"]))
			permissions, _ := adminDAO.GetAdminByAccountID(accountID)
			ctx := context.WithValue(r.Context(), "UserAccessToken", token)
			ctx = context.WithValue(ctx, "id", userInfo["id"])
			ctx = context.WithValue(ctx, "role_name", userInfo["role_name"])
			//set permission into context
			ctx = context.WithValue(ctx, "can_manage_admin", permissions.CanManageAdmin)
			ctx = context.WithValue(ctx, "can_manage_expert", permissions.CanManageExpert)
			ctx = context.WithValue(ctx, "can_manage_moderator", permissions.CanManageModerator)
			ctx = context.WithValue(ctx, "can_manage_learner", permissions.CanManageLearner)
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
			isSuspended := CheckIfSuspended(userInfo["id"])
			if isSuspended {
				http.Error(w, "Your account has been suspended.", http.StatusUnauthorized)
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
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Unauthorized", http.StatusForbidden)
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

func CheckIfSuspended(id interface{}) bool {
	accountDAO := daos.GetAccountDAO()
	accountID, _ := uuid.Parse(fmt.Sprint(id))
	accountInfo, _ := accountDAO.FindAccountByID(accountID)
	if accountInfo.IsSuspended {
		return true
	}
	return false
}
