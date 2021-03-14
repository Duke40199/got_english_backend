package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang/got_english_backend/models"

	"github.com/dgrijalva/jwt-go"
)

//DecodeFirebaseIDToken will decode ID Token from firebase
func DecodeFirebaseIDToken(w http.ResponseWriter, r *http.Request) models.LoginResponse {
	defer r.Body.Close()
	//Split the header to only get the token
	authorizationToken := r.Header.Get("Authorization")
	idToken := strings.TrimSpace(strings.Replace(authorizationToken, "Bearer", "", 1))
	//Parse the token here
	token, _, err := new(jwt.Parser).ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
		return models.LoginResponse{}
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		var loginResponse = models.LoginResponse{
			Username: fmt.Sprintf("%v", claims["username"]),
			RoleName: fmt.Sprintf("%v", claims["role_name"]),
			Iss:      fmt.Sprintf("%v", claims["iss"]),
			Aud:      fmt.Sprintf("%v", claims["aud"]),
		}
		fmt.Print(loginResponse.Username)
		return loginResponse
	}
	return models.LoginResponse{}
}

//DecodeFirebaseIDToken will decode ID Token from firebase
func DecodeFirebaseCustomToken(w http.ResponseWriter, r *http.Request) models.LoginResponse {
	defer r.Body.Close()
	//Split the header to only get the token
	authorizationToken := r.Header.Get("Authorization")
	idToken := strings.TrimSpace(strings.Replace(authorizationToken, "Bearer", "", 1))
	//Parse the token here
	token, _, err := new(jwt.Parser).ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
		return models.LoginResponse{}
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userInfo := claims["claims"].(map[string]interface{})
		var loginResponse = models.LoginResponse{
			Username: fmt.Sprintf("%v", userInfo["username"]),
			RoleName: fmt.Sprintf("%v", userInfo["role_name"]),
			Iss:      fmt.Sprintf("%v", claims["iss"]),
			Aud:      fmt.Sprintf("%v", claims["aud"]),
		}
		fmt.Print(loginResponse.Username)
		return loginResponse
	}
	return models.LoginResponse{}
}
