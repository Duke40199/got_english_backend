package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang/got_english_backend/models"

	"github.com/dgrijalva/jwt-go"
)

//DecodeFirebaseIDToken will decode ID Token from firebase
func DecodeGoogleToken(w http.ResponseWriter, r *http.Request) models.GoogleTokenStruct {
	defer r.Body.Close()
	//Split the header to only get the token
	authorizationToken := r.Header.Get("Authorization")
	idToken := strings.TrimSpace(strings.Replace(authorizationToken, "Bearer", "", 1))
	//Parse the token here
	token, _, err := new(jwt.Parser).ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return models.GoogleTokenStruct{}
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		emailVerified, _ := strconv.ParseBool(fmt.Sprintf("%v", claims["email"]))
		var googleIDTokenStruct = models.GoogleTokenStruct{
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
	return models.GoogleTokenStruct{}
}

//DecodeFirebaseIDToken will decode ID Token from firebase
func DecodeFirebaseCustomToken(w http.ResponseWriter, r *http.Request) models.FirebaseCustomTokenStruct {
	defer r.Body.Close()
	//Split the header to only get the token
	authorizationToken := r.Header.Get("Authorization")
	idToken := strings.TrimSpace(strings.Replace(authorizationToken, "Bearer", "", 1))
	//Parse the token here
	token, _, err := new(jwt.Parser).ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return models.FirebaseCustomTokenStruct{}
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userInfo := claims["claims"].(map[string]interface{})
		var firebaseCustomTokenStruct = models.FirebaseCustomTokenStruct{
			Username: fmt.Sprintf("%v", userInfo["username"]),
			RoleName: fmt.Sprintf("%v", userInfo["role_name"]),
			Iss:      fmt.Sprintf("%v", claims["iss"]),
			Aud:      fmt.Sprintf("%v", claims["aud"]),
		}
		return firebaseCustomTokenStruct
	}
	return models.FirebaseCustomTokenStruct{}
}

//DecodeFirebaseIDToken will decode ID Token from firebase
func DecodeFirebaseIDToken(w http.ResponseWriter, r *http.Request) models.FirebaseIDTokenStruct {
	defer r.Body.Close()
	//Split the header to only get the token
	authorizationToken := r.Header.Get("Authorization")
	idToken := strings.TrimSpace(strings.Replace(authorizationToken, "Bearer", "", 1))
	//Parse the token here
	token, _, err := new(jwt.Parser).ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return models.FirebaseIDTokenStruct{}
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		var firebaseIDTokenStruct = models.FirebaseIDTokenStruct{
			Iss:    fmt.Sprintf("%v", claims["iss"]),
			Aud:    fmt.Sprintf("%v", claims["aud"]),
			Sub:    fmt.Sprintf("%v", claims["sub"]),
			Email:  fmt.Sprintf("%v", claims["email"]),
			UserID: fmt.Sprintf("%v", claims["user_id"]),
		}
		return firebaseIDTokenStruct
	}
	return models.FirebaseIDTokenStruct{}
}
