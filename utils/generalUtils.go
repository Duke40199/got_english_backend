package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang/got_english_backend/models"

	"github.com/dgrijalva/jwt-go"
)

// GetCurrentTime function is used to get the current time in milliseconds.
func GetCurrentEpochTimeInMiliseconds() int64 {
	var now = time.Now()
	ts := now.UnixNano() / 1000000
	return ts
}

//DecodeFirebaseIDToken will decode ID Token from firebase
func DecodeGoogleIDToken(w http.ResponseWriter, r *http.Request) models.GoogleIDTokenStruct {
	defer r.Body.Close()
	//Split the header to only get the token
	authorizationToken := r.Header.Get("Authorization")
	idToken := strings.TrimSpace(strings.Replace(authorizationToken, "Bearer", "", 1))
	//Parse the token here
	token, _, err := new(jwt.Parser).ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return models.GoogleIDTokenStruct{}
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		emailVerified, _ := strconv.ParseBool(fmt.Sprintf("%v", claims["email"]))
		var googleIDTokenStruct = models.GoogleIDTokenStruct{
			Iss:           fmt.Sprintf("%v", claims["iss"]),
			Azp:           fmt.Sprintf("%v", claims["azp"]),
			Aud:           fmt.Sprintf("%v", claims["aud"]),
			Sub:           fmt.Sprintf("%v", claims["sub"]),
			Email:         fmt.Sprintf("%v", claims["email"]),
			EmailVerified: emailVerified,
			Name:          fmt.Sprintf("%v", claims["name"]),
			Picture:       fmt.Sprintf("%v", claims["picture"]),
			GivenName:     fmt.Sprintf("%v", claims["given_name"]),
			FamilyName:    fmt.Sprintf("%v", claims["family_name"]),
			Locale:        fmt.Sprintf("%v", claims["locale"]),
		}
		return googleIDTokenStruct
	}
	return models.GoogleIDTokenStruct{}
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
		return loginResponse
	}
	return models.LoginResponse{}
}
