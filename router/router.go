package router

import (
	"log"
	"net/http"

	"github.com/golang/got_english_backend/middleware"

	"github.com/golang/got_english_backend/config"
	"github.com/golang/got_english_backend/controllers"

	"github.com/gorilla/mux"
)

var r = mux.NewRouter()

func init() {
	log.Println("Initializing Router")
	apiV1 := r.PathPrefix("/").Subrouter()
	//For root functions
	apiV1.HandleFunc("/", RootRoute).Methods("GET")
	apiV1.HandleFunc("/login", controllers.LoginHandler).Methods("POST")
	apiV1.HandleFunc("/profile", middleware.FirebaseAuthentication(controllers.ViewProfileHandler)).Methods("GET")
	//For user functions
	apiV1.HandleFunc("/users", middleware.FirebaseAuthentication(controllers.ViewProfileHandler)).Methods("GET")
	apiV1.HandleFunc("/users/register", middleware.FirebaseAuthentication(controllers.CreateUserHandler)).Methods("POST")
	// For changelogs
	//apiV1.HandleFunc("/cms/changelogs", controller.ShowChangelogsHandler).Methods("GET")
	//apiV1.HandleFunc("/cms/health-check", controller.GetHealthCHeck).Methods("GET")
	//// For swagger
	//apiV1.HandleFunc("/cms/swagger", middleware.ProductAuthentication(controller.GetSwaggerHandler)).Methods("GET")
	log.Println("Initialed Router")
}

func GetRouter() *mux.Router {
	return r
}

func RootRoute(w http.ResponseWriter, r *http.Request) {
	config.ResponseWithSuccess(w, "OK", "Welcome to this API.")
}
