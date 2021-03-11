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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf(`
	-----------------------------------------------------
	App name: %s
	Version: %s
	Listening Port: %v
	Environment: %s
	-----------------------------------------------------
	`, config.AppName, config.AppVersion, port, config.Environment)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	})

	log.Fatal(http.ListenAndServe(":"+port, c.Handler(r)))

}

func init() {
}
