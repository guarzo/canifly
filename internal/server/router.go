package server

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	flyHandlers "github.com/guarzo/canifly/internal/handlers"
	flyHttp "github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

// SetupHandlers configures and returns the app's router
func SetupHandlers(secret string, logger interfaces.Logger, appServices *AppServices, basePath string) http.Handler {
	sessionStore := flyHttp.NewSessionService(secret)
	r := mux.NewRouter()

	// Add authentication middleware
	r.Use(flyHttp.AuthMiddleware(sessionStore, logger))

	authHandler := flyHandlers.NewAuthHandler(sessionStore, appServices.EVEDataService, logger, appServices.AccountManagementService, appServices.ConfigurationService, appServices.LoginService, appServices.AuthClient)
	accountHandler := flyHandlers.NewAccountHandler(sessionStore, logger, appServices.AccountManagementService, appServices.HTTPCacheService, appServices.WebSocketHub)
	characterHandler := flyHandlers.NewCharacterHandler(logger, appServices.EVEDataService, appServices.HTTPCacheService)
	skillPlanHandler := flyHandlers.NewSkillPlanHandler(logger, appServices.EVEDataService, appServices.AccountManagementService, appServices.HTTPCacheService, appServices.WebSocketHub)
	configHandler := flyHandlers.NewConfigHandler(logger, appServices.ConfigurationService, appServices.HTTPCacheService)
	eveDataHandler := flyHandlers.NewEveDataHandler(logger, appServices.SyncService, appServices.ConfigurationService, appServices.EVEDataService, appServices.AccountManagementService, appServices.HTTPCacheService)
	assocHandler := flyHandlers.NewAssociationHandler(logger, appServices.AccountManagementService)
	fuzzworksHandler := flyHandlers.NewFuzzworksHandler(logger, basePath, appServices.HTTPCacheService, appServices.WebSocketHub)

	// Public routes
	r.HandleFunc("/callback", authHandler.CallBack())
	r.HandleFunc("/api/add-character", authHandler.AddCharacterHandler())
	r.HandleFunc("/api/finalize-login", authHandler.FinalizeLogin())

	// Auth routes

	r.HandleFunc("/api/session", authHandler.GetSession()).Methods("GET")
	r.HandleFunc("/api/logout", authHandler.Logout())
	r.HandleFunc("/api/login", authHandler.Login())
	r.HandleFunc("/api/reset-identities", authHandler.ResetAccounts())

	// RESTful skill plan endpoints
	r.HandleFunc("/api/skill-plans", skillPlanHandler.ListSkillPlans()).Methods("GET")
	r.HandleFunc("/api/skill-plans", skillPlanHandler.CreateSkillPlan()).Methods("POST")
	r.HandleFunc("/api/skill-plans/refresh", skillPlanHandler.RefreshSkillPlans()).Methods("POST")
	r.HandleFunc("/api/skill-plans/{name}", skillPlanHandler.GetSkillPlan()).Methods("GET")
	r.HandleFunc("/api/skill-plans/{name}", skillPlanHandler.UpdateSkillPlan()).Methods("PUT")
	r.HandleFunc("/api/skill-plans/{name}", skillPlanHandler.DeleteSkillPlanRESTful()).Methods("DELETE")
	r.HandleFunc("/api/skill-plans/{name}/copy", skillPlanHandler.CopySkillPlan()).Methods("POST")

	// RESTful account endpoints
	r.HandleFunc("/api/accounts", accountHandler.ListAccounts()).Methods("GET")
	r.HandleFunc("/api/accounts/{id}", accountHandler.GetAccount()).Methods("GET")
	r.HandleFunc("/api/accounts/{id}", accountHandler.UpdateAccount()).Methods("PATCH")
	r.HandleFunc("/api/accounts/{id}", accountHandler.DeleteAccount()).Methods("DELETE")

	// RESTful character endpoints
	r.HandleFunc("/api/characters/{id}", characterHandler.GetCharacter()).Methods("GET")
	r.HandleFunc("/api/characters/{id}", characterHandler.UpdateCharacterRESTful()).Methods("PATCH")
	r.HandleFunc("/api/characters/{id}", characterHandler.DeleteCharacter()).Methods("DELETE")
	r.HandleFunc("/api/characters/{id}/refresh", characterHandler.RefreshCharacter()).Methods("POST")

	// RESTful config endpoints
	r.HandleFunc("/api/config", configHandler.GetConfig()).Methods("GET")
	r.HandleFunc("/api/config", configHandler.UpdateConfig()).Methods("PATCH")
	r.HandleFunc("/api/config/eve/status", configHandler.GetEVEConfigStatus()).Methods("GET")
	r.HandleFunc("/api/config/eve/credentials", configHandler.SaveEVECredentials()).Methods("POST")

	// EVE data endpoints
	r.HandleFunc("/api/eve/skill-plans", eveDataHandler.GetSkillPlans).Methods("GET")
	r.HandleFunc("/api/eve/profiles", eveDataHandler.GetEveProfiles).Methods("GET")
	r.HandleFunc("/api/eve/conversions", eveDataHandler.GetEveConversions).Methods("GET")

	r.HandleFunc("/api/sync-subdirectory", eveDataHandler.SyncSubDirectory)
	r.HandleFunc("/api/sync-all-subdirectories", eveDataHandler.SyncAllSubdirectories)
	r.HandleFunc("/api/backup-directory", eveDataHandler.BackupDirectory)

	// RESTful association endpoints
	r.HandleFunc("/api/associations", assocHandler.ListAssociations()).Methods("GET")
	r.HandleFunc("/api/associations", assocHandler.CreateAssociation()).Methods("POST")
	r.HandleFunc("/api/associations/{userId}/{characterId}", assocHandler.DeleteAssociation()).Methods("DELETE")

	// WebSocket endpoint
	r.HandleFunc("/api/ws", appServices.WebSocketHub.HandleWebSocket)

	// Fuzzworks endpoints
	r.HandleFunc("/api/fuzzworks/update", fuzzworksHandler.UpdateData()).Methods("POST")
	r.HandleFunc("/api/fuzzworks/status", fuzzworksHandler.GetStatus()).Methods("GET")

	// Serve static files from filesystem
	staticDir := "./static"
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		// Try relative to executable location
		if ex, err := os.Executable(); err == nil {
			staticDir = filepath.Join(filepath.Dir(ex), "static")
		}
	}
	staticFileServer := http.FileServer(http.Dir(staticDir))
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
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With"}),
		handlers.AllowCredentials(),
	)(h)
}
