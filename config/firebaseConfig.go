package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

//FirebaseConfig struct
type FirebaseConfig struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}

//GenerateFirebaseConfigJSON from env
func GenerateFirebaseConfigJSON() []byte {
	config := FirebaseConfig{
		Type:                    os.Getenv("FIREBASE_TYPE"),
		ProjectID:               os.Getenv("FIREBASE_PROJECT_ID"),
		PrivateKeyID:            os.Getenv("FIREBASE_PRIVATE_KEY_ID"),
		PrivateKey:              os.Getenv("FIREBASE_PRIVATE_KEY"),
		ClientEmail:             os.Getenv("FIREBASE_CLIENT_EMAIL"),
		ClientID:                os.Getenv("FIREBASE_CLIENT_ID"),
		AuthURI:                 os.Getenv("FIREBASE_AUTH_URI"),
		TokenURI:                os.Getenv("FIREBASE_TOKEN_URI"),
		AuthProviderX509CertURL: os.Getenv("FIREBASE_AUTH_PROVIDER_X509_CERT_URL"),
		ClientX509CertURL:       os.Getenv("FIREBASE_CLIENT_X509_CERT_URL"),
	}
	jsonByte, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("======CONFIG: %s", jsonByte)
	return jsonByte
}

//SetupFirebase based on env config
func SetupFirebase() (*auth.Client, context.Context) {
	var context = context.Background()
	// jsonByte := GenerateFirebaseConfigJSON()
	opt := option.WithCredentialsFile("gotenglish-app-firebase-adminsdk-hmhes-82eb0253df.json")
	// opt := option.WithCredentialsJSON(GenerateFirebaseConfigJSON())
	//Firebase admin SDK initialization
	app, err := firebase.NewApp(context, nil, opt)
	if err != nil {
		panic(err)
	}

	//Firebase Auth
	client, err := app.Auth(context)
	if err != nil {
		log.Fatalf("======== ERROR:%s", err)
	}
	fmt.Println("FIREBASE COMPLETE")
	return client, context
}
