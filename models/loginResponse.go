package models

import "github.com/dgrijalva/jwt-go"

//LoginResponse struct
type LoginResponse struct {
	Foo string `json:"foo"`
	jwt.StandardClaims
	RoleName string `json:"role_name"`
	Username string `json:"username"`
	Iss      string `json:"iss"`
	Aud      string `json:"aud"`
	AuthTime int    `json:"auth_time"`
	UserID   string `json:"user_id"`
	Sub      string `json:"sub"`
	Iat      int    `json:"iat"`
	Exp      int    `json:"exp"`
	Firebase struct {
		Identities struct {
		} `json:"identities"`
		SignInProvider string `json:"sign_in_provider"`
	} `json:"firebase"`
}
