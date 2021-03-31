package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/got_english_backend/config"
)

var roleNameConfig = config.GetRoleNameConfig()

// FirebaseAuthentication : to verify all authorized operations
func UserAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			authorizationToken := r.Header.Get("Authorization")
			customToken := strings.TrimSpace(strings.Replace(authorizationToken, "Bearer", "", 1))
			token, _ := jwt.Parse(customToken, nil)
			if token == nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			claims, _ := token.Claims.(jwt.MapClaims)
			userInfo := claims["claims"].(map[string]interface{})
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
			authorizationToken := r.Header.Get("Authorization")
			customToken := strings.TrimSpace(strings.Replace(authorizationToken, "Bearer", "", 1))
			token, _ := jwt.Parse(customToken, nil)
			if token == nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			claims, _ := token.Claims.(jwt.MapClaims)
			userInfo := claims["claims"].(map[string]interface{})
			if userInfo["role_name"] != roleNameConfig.Moderator {
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

func LearnerAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			authorizationToken := r.Header.Get("Authorization")
			customToken := strings.TrimSpace(strings.Replace(authorizationToken, "Bearer", "", 1))
			token, _ := jwt.Parse(customToken, nil)
			if token == nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			claims, _ := token.Claims.(jwt.MapClaims)
			userInfo := claims["claims"].(map[string]interface{})
			if userInfo["role_name"] != roleNameConfig.Learner {
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
				return
			}
			claims, _ := token.Claims.(jwt.MapClaims)
			userInfo := claims["claims"].(map[string]interface{})
			if userInfo["role_name"] != roleNameConfig.Admin {
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

// FirebaseAuthentication : to verify all authorized operations
func ExpertAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			authorizationToken := r.Header.Get("Authorization")
			customToken := strings.TrimSpace(strings.Replace(authorizationToken, "Bearer", "", 1))
			//parse
			token, _ := jwt.Parse(customToken, nil)
			if token == nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			claims, _ := token.Claims.(jwt.MapClaims)
			userInfo := claims["claims"].(map[string]interface{})
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
