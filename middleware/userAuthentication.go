package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/utils"
)

var roleNameConfig = config.GetRoleNameConfig()

// FirebaseAuthentication : to verify all authorized operations
func FirebaseAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorizationToken := r.Header.Get("Authorization")
		firebaseAuth, _ := config.SetupFirebase()
		idToken := strings.TrimSpace(strings.Replace(authorizationToken, "Bearer", "", 1))
		if idToken == "" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		//verify token
		token, err := firebaseAuth.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			log.Printf("====ERR: %s", err)
			http.Error(w, "Invalid Token!", http.StatusForbidden)
			return
		}
		//verify roleName
		loginResponse := utils.DecodeFirebaseIDToken(w, r)
		if loginResponse.RoleName == "" {
			http.Error(w, "Your current role cannot access this function.", http.StatusForbidden)
		} else {
			ctx := context.WithValue(r.Context(), "UserAccessToken", token)
			ctx = context.WithValue(ctx, "UUID", token.UID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

// FirebaseAuthentication : to verify all authorized operations
func UserAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			authorizationToken := r.Header.Get("Authorization")
			customToken := strings.TrimSpace(strings.Replace(authorizationToken, "Bearer", "", 1))
			//parse
			token, _ := jwt.Parse(customToken, nil)
			if token == nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
			} else {
				claims, _ := token.Claims.(jwt.MapClaims)
				userInfo := claims["claims"].(map[string]interface{})
				if userInfo["role_name"] == "" {
					http.Error(w, "Your current role cannot access this function.", http.StatusForbidden)
				} else {
					ctx := context.WithValue(r.Context(), "UserAccessToken", token)
					ctx = context.WithValue(ctx, "UUID", userInfo["username"])
					next.ServeHTTP(w, r.WithContext(ctx))
				}
			}
		} else {
			http.Error(w, "Unauthorized", http.StatusForbidden)
		}
	}
}

func ModeratorAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			authorizationToken := r.Header.Get("Authorization")
			customToken := strings.TrimSpace(strings.Replace(authorizationToken, "Bearer", "", 1))
			//parse
			token, _ := jwt.Parse(customToken, nil)
			if token == nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
			} else {
				claims, _ := token.Claims.(jwt.MapClaims)
				userInfo := claims["claims"].(map[string]interface{})
				if userInfo["role_name"] == roleNameConfig.Moderator {
					http.Error(w, "Your current role cannot access this function.", http.StatusForbidden)
				} else {
					ctx := context.WithValue(r.Context(), "UserAccessToken", token)
					ctx = context.WithValue(ctx, "UUID", userInfo["username"])
					next.ServeHTTP(w, r.WithContext(ctx))
				}
			}
		} else {
			http.Error(w, "Unauthorized", http.StatusForbidden)
		}
	}
}

// FirebaseAuthentication : to verify all authorized operations
func AdminAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			authorizationToken := r.Header.Get("Authorization")
			customToken := strings.TrimSpace(strings.Replace(authorizationToken, "Bearer", "", 1))
			//parse
			token, _ := jwt.Parse(customToken, nil)
			if token == nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
			} else {
				claims, _ := token.Claims.(jwt.MapClaims)
				userInfo := claims["claims"].(map[string]interface{})
				if userInfo["role_name"] == roleNameConfig.Admin {
					http.Error(w, "Your current role cannot access this function.", http.StatusForbidden)
				} else {
					ctx := context.WithValue(r.Context(), "UserAccessToken", token)
					ctx = context.WithValue(ctx, "UUID", userInfo["username"])
					next.ServeHTTP(w, r.WithContext(ctx))
				}
			}
		} else {
			http.Error(w, "Unauthorized", http.StatusForbidden)
		}
	}
}
