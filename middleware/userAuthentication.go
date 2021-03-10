package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/utils"
)

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
		}
		ctx := context.WithValue(r.Context(), "UserAccessToken", token)
		ctx = context.WithValue(ctx, "UUID", token.UID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

}
