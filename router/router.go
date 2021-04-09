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
	//For administrator (web admin) functions
	apiV1.HandleFunc("/administrator/summary", middleware.UserAuthentication(controllers.GetAdministratorSummary)).Methods("GET")

	//For root functions
	apiV1.HandleFunc("/", RootRoute).Methods("GET")
	apiV1.HandleFunc("/login", controllers.LoginHandler).Methods("POST")
	apiV1.HandleFunc("/login/google", controllers.LoginWithGoogleHandler).Methods("POST")
	apiV1.HandleFunc("/profile", middleware.UserAuthentication(controllers.ViewProfileHandler)).Methods("GET")

	//For account functions
	apiV1.HandleFunc("/accounts", middleware.AdminAuthentication(controllers.GetAccountsHandler)).Methods("GET")
	apiV1.HandleFunc("/accounts", controllers.CreateAccountHandler).Methods("POST")
	apiV1.HandleFunc("/accounts/{account_id}/update", middleware.UserAuthentication(controllers.UpdateAccountHandler)).Methods("PUT")

	//For admin functions
	apiV1.HandleFunc("/admins/{account_id}/update", middleware.AdminAuthentication(controllers.UpdateAdminHandler)).Methods("PUT")
	//For application form functions
	apiV1.HandleFunc("/application-forms", middleware.ExpertAuthentication(controllers.CreateApplicationFormHandler)).Methods("POST")
	apiV1.HandleFunc("/application-forms", middleware.ModeratorAuthentication(controllers.GetApplicationFormsHandler)).Methods("GET")
	apiV1.HandleFunc("/application-forms/history", middleware.ExpertAuthentication(controllers.GetApplicationFormsHandler)).Methods("GET")

	//For coin bundle functions
	apiV1.HandleFunc("/coin-bundles", middleware.UserAuthentication(controllers.GetCoinBundlesHandler)).Methods("GET")
	apiV1.HandleFunc("/coin-bundles", middleware.ModeratorAuthentication(controllers.CreateCoinBundleHandler)).Methods("POST")
	apiV1.HandleFunc("/coin-bundles/{coin_bundle_id}/update", middleware.ModeratorAuthentication(controllers.UpdateCoinBundleHandler)).Methods("PUT")

	//For expert functions
	apiV1.HandleFunc("/experts/{account_id}/update", middleware.AdminAuthentication(controllers.UpdateExpertHandler)).Methods("PUT")
	apiV1.HandleFunc("/experts/translators", middleware.LearnerAuthentication(controllers.GetTranslatorExpertsHandler)).Methods("GET")
	//For invoice functions
	apiV1.HandleFunc("/invoices", middleware.LearnerAuthentication(controllers.CreateInvoiceHandler)).Methods("POST")
	apiV1.HandleFunc("/invoices/history", middleware.LearnerAuthentication(controllers.CreateMessagingSessionHandler)).Methods("GET")
	apiV1.HandleFunc("/invoices/{invoice_id}", middleware.LearnerAuthentication(controllers.CreateMessagingSessionHandler)).Methods("GET")
	apiV1.HandleFunc("/invoices/{invoice_id}/update", middleware.LearnerExpertAuthentication(controllers.UpdateMessagingSessionHandler)).Methods("PUT")

	//For messaging session functions
	apiV1.HandleFunc("/messaging-sessions", middleware.UserAuthentication(controllers.GetMessagingSessionHandler)).Methods("GET")
	apiV1.HandleFunc("/messaging-sessions", middleware.LearnerAuthentication(controllers.CreateMessagingSessionHandler)).Methods("POST")
	apiV1.HandleFunc("/messaging-sessions/{messaging_session_id}/update", middleware.LearnerExpertAuthentication(controllers.UpdateMessagingSessionHandler)).Methods("PUT")
	apiV1.HandleFunc("/messaging-sessions/{messaging_session_id}/cancel", middleware.LearnerAuthentication(controllers.CancelMessagingSessionHandler)).Methods("PUT")
	//For moderator functions
	apiV1.HandleFunc("/moderators/{account_id}/update", middleware.AdminAuthentication(controllers.UpdateModeratorHandler)).Methods("PUT")

	//For pricing functions
	apiV1.HandleFunc("/pricings", middleware.ModeratorAuthentication(controllers.GetPricingsHandler)).Methods("GET")
	apiV1.HandleFunc("/pricings", middleware.ModeratorAuthentication(controllers.CreatePricingHandler)).Methods("POST")
	apiV1.HandleFunc("/pricings/{pricing_id}/update", middleware.ModeratorAuthentication(controllers.UpdatePricingHandler)).Methods("PUT")
	apiV1.HandleFunc("/pricings/{pricing_id}/delete", middleware.ModeratorAuthentication(controllers.DeletePricingHandler)).Methods("DELETE")
	//For live call session functions
	apiV1.HandleFunc("/live-call-sessions", middleware.UserAuthentication(controllers.GetLiveCallSessionsHandler)).Methods("GET")
	apiV1.HandleFunc("/live-call-sessions", middleware.LearnerAuthentication(controllers.CreateLiveCallSessionHandler)).Methods("POST")
	apiV1.HandleFunc("/live-call-sessions/{live_call_session_id}/update", middleware.LearnerExpertAuthentication(controllers.UpdateLiveCallSessionHandler)).Methods("PUT")
	apiV1.HandleFunc("/live-call-sessions/{live_call_session_id}/cancel", middleware.LearnerAuthentication(controllers.CancelLiveCallHandler)).Methods("PUT")

	//For translation session functions
	// apiV1.HandleFunc("/translation-sessions", middleware.UserAuthentication(controllers.GetLiveCallSessionsHandler)).Methods("GET")
	apiV1.HandleFunc("/translation-sessions", middleware.LearnerAuthentication(controllers.CreateTranslationSessionHandler)).Methods("POST")
	apiV1.HandleFunc("/translation-sessions/{translation_session_id}/update", middleware.LearnerExpertAuthentication(controllers.UpdateTranslationSessionHandler)).Methods("PUT")
	apiV1.HandleFunc("/translation-sessions/{translation_session_id}/cancel", middleware.LearnerAuthentication(controllers.CancelTranslationSessionHandler)).Methods("PUT")
	//For ratings functions
	apiV1.HandleFunc("/ratings", middleware.AdminAuthentication(controllers.GetRatingsHandler)).Methods("GET")
	apiV1.HandleFunc("/ratings", middleware.LearnerAuthentication(controllers.CreateRatingHandler)).Methods("POST")

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
