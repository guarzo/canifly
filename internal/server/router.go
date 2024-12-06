package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/guarzo/canifly/internal/embed"
	flyHandlers "github.com/guarzo/canifly/internal/handlers"
	http2 "github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/services"
	"github.com/guarzo/canifly/internal/services/esi"
)

// SetupRouter configures and returns the app’s router
func SetupRouter(secret string, logger *logrus.Logger, esiService esi.ESIService, skillService *services.SkillService, configService *services.ConfigService) *mux.Router {
	sessionStore := http2.NewSessionService(secret)
	r := mux.NewRouter()

	// Add authentication middleware
	r.Use(http2.AuthMiddleware(sessionStore, logger))
	dashboardHandler := flyHandlers.NewDashboardHandler(sessionStore, esiService, logger, skillService)
	authHandler := flyHandlers.NewAuthHandler(sessionStore, esiService, logger)
	accountHandler := flyHandlers.NewAccountHandler(sessionStore, logger)
	characterHandler := flyHandlers.NewCharacterHandler(sessionStore, logger, configService)
	skillPlanHandler := flyHandlers.NewSkillPlanHandler(logger, skillService)

	// Define routes
	r.HandleFunc("/callback/", authHandler.CallBack())
	r.HandleFunc("/login", authHandler.Login())
	r.HandleFunc("/logout", authHandler.Logout())
	r.HandleFunc("/auth-character", authHandler.AddCharacterHandler())
	r.HandleFunc("/api/home-data", dashboardHandler.GetDashboardData()).Methods("GET")
	r.HandleFunc("/api/logout", authHandler.Logout())
	r.HandleFunc("/api/get-skill-plan", skillPlanHandler.GetSkillPlanFile())
	r.HandleFunc("/api/save-skill-plan", skillPlanHandler.SaveSkillPlan())
	r.HandleFunc("/api/delete-skill-plan", skillPlanHandler.DeleteSkillPlan())
	r.HandleFunc("/api/unassigned-characters", characterHandler.GetUnassignedCharacters())
	r.HandleFunc("/api/assign-character", characterHandler.AssignCharacter())
	r.HandleFunc("/api/update-account-name", accountHandler.UpdateAccountName())
	r.HandleFunc("/api/update-character", characterHandler.UpdateCharacter())
	r.HandleFunc("/api/toggle-account-status", accountHandler.ToggleAccountStatus())
	r.HandleFunc("/api/remove-character", characterHandler.RemoveCharacter())
	r.HandleFunc("/reset-identities", authHandler.ResetAccounts())

	// Serve static files
	staticFileServer := http.FileServer(http.FS(embed.StaticFilesSub))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticFileServer))

	return r
}
