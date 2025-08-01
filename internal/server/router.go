package server

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/guarzo/canifly/internal/embed"
	flyHandlers "github.com/guarzo/canifly/internal/handlers"
	flyHttp "github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

// SetupHandlers configures and returns the app’s router
func SetupHandlers(secret string, logger interfaces.Logger, appServices *AppServices) http.Handler {
	sessionStore := flyHttp.NewSessionService(secret)
	r := mux.NewRouter()

	// Add authentication middleware
	r.Use(flyHttp.AuthMiddleware(sessionStore, logger))

	dashboardHandler := flyHandlers.NewDashboardHandler(sessionStore, logger, appServices.ConfigurationService)
	authHandler := flyHandlers.NewAuthHandler(sessionStore, appServices.EVEDataService, logger, appServices.AccountManagementService, appServices.ConfigurationService, appServices.LoginService, appServices.AuthClient)
	accountHandler := flyHandlers.NewAccountHandler(sessionStore, logger, appServices.AccountManagementService)
	characterHandler := flyHandlers.NewCharacterHandler(logger, appServices.EVEDataService)
	skillPlanHandler := flyHandlers.NewSkillPlanHandler(logger, appServices.EVEDataService)
	configHandler := flyHandlers.NewConfigHandler(logger, appServices.ConfigurationService)
	eveDataHandler := flyHandlers.NewEveDataHandler(logger, appServices.SyncService)
	assocHandler := flyHandlers.NewAssociationHandler(logger, appServices.AccountManagementService)

	// Public routes
	r.HandleFunc("/callback/", authHandler.CallBack())
	r.HandleFunc("/api/add-character", authHandler.AddCharacterHandler())
	r.HandleFunc("/api/finalize-login", authHandler.FinalizeLogin())

	// Auth routes
	r.HandleFunc("/api/app-data", dashboardHandler.GetDashboardData()).Methods("GET")
	r.HandleFunc("/api/app-data-no-cache", dashboardHandler.GetDashboardDataNoCache()).Methods("GET")

	r.HandleFunc("/api/logout", authHandler.Logout())
	r.HandleFunc("/api/login", authHandler.Login())
	r.HandleFunc("/api/reset-identities", authHandler.ResetAccounts())

	r.HandleFunc("/api/get-skill-plan", skillPlanHandler.GetSkillPlanFile())
	r.HandleFunc("/api/save-skill-plan", skillPlanHandler.SaveSkillPlan())
	r.HandleFunc("/api/delete-skill-plan", skillPlanHandler.DeleteSkillPlan())

	// RESTful account endpoints
	r.HandleFunc("/api/accounts", accountHandler.ListAccounts()).Methods("GET")
	r.HandleFunc("/api/accounts/{id}", accountHandler.GetAccount()).Methods("GET")
	r.HandleFunc("/api/accounts/{id}", accountHandler.UpdateAccount()).Methods("PATCH")
	r.HandleFunc("/api/accounts/{id}", accountHandler.DeleteAccount()).Methods("DELETE")
	
	// Legacy account endpoints (to be deprecated)
	r.HandleFunc("/api/update-account-name", accountHandler.UpdateAccountName())
	r.HandleFunc("/api/toggle-account-status", accountHandler.ToggleAccountStatus())
	r.HandleFunc("/api/toggle-account-visibility", accountHandler.ToggleAccountVisibility())
	r.HandleFunc("/api/remove-account", accountHandler.RemoveAccount())

	// RESTful character endpoints
	r.HandleFunc("/api/characters/{id}", characterHandler.GetCharacter()).Methods("GET")
	r.HandleFunc("/api/characters/{id}", characterHandler.UpdateCharacterRESTful()).Methods("PATCH")
	r.HandleFunc("/api/characters/{id}", characterHandler.DeleteCharacter()).Methods("DELETE")
	
	// Legacy character endpoints (to be deprecated)
	r.HandleFunc("/api/update-character", characterHandler.UpdateCharacter)
	r.HandleFunc("/api/remove-character", characterHandler.RemoveCharacter)

	// RESTful config endpoints
	r.HandleFunc("/api/config", configHandler.GetConfig()).Methods("GET")
	r.HandleFunc("/api/config", configHandler.UpdateConfig()).Methods("PATCH")
	
	// Legacy config endpoints (to be deprecated)
	r.HandleFunc("/api/choose-settings-dir", configHandler.ChooseSettingsDir)
	r.HandleFunc("/api/reset-to-default-directory", configHandler.ResetToDefaultDir)
	r.HandleFunc("/api/save-user-selections", configHandler.SaveUserSelections)

	r.HandleFunc("/api/sync-subdirectory", eveDataHandler.SyncSubDirectory)
	r.HandleFunc("/api/sync-all-subdirectories", eveDataHandler.SyncAllSubdirectories)
	r.HandleFunc("/api/backup-directory", eveDataHandler.BackupDirectory)

	r.HandleFunc("/api/associate-character", assocHandler.AssociateCharacter)
	r.HandleFunc("/api/unassociate-character", assocHandler.UnassociateCharacter)

	// Serve static files
	staticFileServer := http.FileServer(http.FS(embed.StaticFilesSub))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticFileServer))

	return createCORSHandler(r)
}

func createCORSHandler(h http.Handler) http.Handler {
	frontendPort := os.Getenv("FRONTEND_PORT")
	if frontendPort == "" {
		frontendPort = "3113" // Default to 3113 if not set
	}
	allowedOrigin := "http://localhost:" + frontendPort
	
	return handlers.CORS(
		handlers.AllowedOrigins([]string{allowedOrigin}),
		handlers.AllowedMethods([]string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With"}),
		handlers.AllowCredentials(),
	)(h)
}
