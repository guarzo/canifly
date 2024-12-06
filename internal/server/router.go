package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/guarzo/canifly/internal/embed"
	flyHandlers "github.com/guarzo/canifly/internal/handlers"
)

// SetupRouter configures and returns the app’s router
func SetupRouter(secret string, logger *logrus.Logger) *mux.Router {
	sessionStore := flyHandlers.NewSessionService(secret)
	r := mux.NewRouter()

	// Add authentication middleware
	r.Use(flyHandlers.AuthMiddleware(sessionStore, logger))

	// Define routes
	r.HandleFunc("/callback/", flyHandlers.CallbackHandler(sessionStore))
	r.HandleFunc("/login", flyHandlers.LoginHandler)
	r.HandleFunc("/auth-character", flyHandlers.AuthCharacterHandler)
	r.HandleFunc("/api/home-data", flyHandlers.HomeDataHandler(sessionStore))
	r.HandleFunc("/api/logout", flyHandlers.LogoutHandler(sessionStore))
	r.HandleFunc("/api/get-skill-plan", flyHandlers.SkillPlanFileHandler)
	r.HandleFunc("/api/save-skill-plan", flyHandlers.SaveSkillPlanHandler)
	r.HandleFunc("/api/delete-skill-plan", flyHandlers.DeleteSkillPlanHandler)
	r.HandleFunc("/api/unassigned-characters", flyHandlers.GetUnassignedCharactersHandler(sessionStore))
	r.HandleFunc("/api/assign-character", flyHandlers.AssignCharacterHandler(sessionStore))
	r.HandleFunc("/api/update-account-name", flyHandlers.UpdateAccountNameHandler(sessionStore))
	r.HandleFunc("/api/update-character", flyHandlers.UpdateCharacterHandler(sessionStore))
	r.HandleFunc("/api/toggle-account-status", flyHandlers.ToggleAccountStatusHandler(sessionStore))
	r.HandleFunc("/api/remove-character", flyHandlers.RemoveCharacterHandler(sessionStore))
	r.HandleFunc("/reset-identities", flyHandlers.ResetAccountsHandler(sessionStore))

	// Serve static files
	staticFileServer := http.FileServer(http.FS(embed.StaticFilesSub))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticFileServer))

	return r
}
