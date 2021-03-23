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
	apiV1.HandleFunc("/profile", middleware.UserAuthentication(controllers.ViewProfileHandler)).Methods("GET")

	//For account functions
	apiV1.HandleFunc("/accounts", middleware.AdminAuthentication(controllers.GetAccountsHandler)).Methods("GET")
	apiV1.HandleFunc("/accounts", controllers.CreateAccountHandler).Methods("POST")
	apiV1.HandleFunc("/accounts/{account_id}/update", middleware.UserAuthentication(controllers.UpdateAccountHandler)).Methods("PUT")
	//For coin bundle functions
	apiV1.HandleFunc("/coin-bundles", middleware.UserAuthentication(controllers.GetCoinBundlesHandler)).Methods("GET")
	apiV1.HandleFunc("/coin-bundles", middleware.ModeratorAuthentication(controllers.CreateCoinBundleHandler)).Methods("POST")
	apiV1.HandleFunc("/coin-bundles/{coin_bundle_id}/update", middleware.ModeratorAuthentication(controllers.UpdateCoinBundleHandler)).Methods("PUT")
	// For changelogs
	//apiV1.HandleFunc("/cms/changelogs", controller.ShowChangelogsHandler).Methods("GET")
	//apiV1.HandleFunc("/cms/health-check", controller.GetHealthCHeck).Methods("GET")
	//// For swagger
	//apiV1.HandleFunc("/cms/swagger", middleware.ProductAuthentication(controller.GetSwaggerHandler)).Methods("GET")
	sh := http.StripPrefix("/swagger/", http.FileServer(http.Dir("./swagger/")))
	r.PathPrefix("/swagger/").Handler(sh)
}

func GetRouter() *mux.Router {
	return r
}

func RootRoute(w http.ResponseWriter, r *http.Request) {
	config.ResponseWithSuccess(w, "OK", "Welcome to this API.")
}
