package main

import (
	"log"
	"net/http"
	"os"

	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/database"
	"github.com/golang/got_english_backend/router"

	"github.com/rs/cors"
)

func main() {
	database.SyncDB(true)
	config := config.GetConfig()
	r := router.GetRouter()

	defaultListeningPort := "4000"
	apiPort := defaultListeningPort

	if port := os.Getenv("PORT"); port != "" {
		apiPort = port
	}

	log.Printf(`
	-----------------------------------------------------
	App name: %s
	Version: %s
	Listening Port: %v
	Environment: %s
	-----------------------------------------------------
	`, config.AppName, config.AppVersion, apiPort, config.Environment)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	})

	log.Fatal(http.ListenAndServe(":"+apiPort, c.Handler(r)))

}

func init() {
}
