package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/guarzo/canifly/internal/embed"
	flyHandlers "github.com/guarzo/canifly/internal/handlers"
	flyHttp "github.com/guarzo/canifly/internal/http"
)

// SetupRouter configures and returns the app’s router
func SetupRouter(secret string, logger *logrus.Logger, appServices *AppServices) *mux.Router {
	sessionStore := flyHttp.NewSessionService(secret)
	r := mux.NewRouter()

	// Add authentication middleware
	r.Use(flyHttp.AuthMiddleware(sessionStore, logger))
	dashboardHandler := flyHandlers.NewDashboardHandler(sessionStore, appServices.EsiService, logger, appServices.SkillService, appServices.DataStore, appServices.ConfigService)
	authHandler := flyHandlers.NewAuthHandler(sessionStore, appServices.EsiService, logger, appServices.DataStore, appServices.ConfigService)
	accountHandler := flyHandlers.NewAccountHandler(sessionStore, logger, appServices.DataStore)
	characterHandler := flyHandlers.NewCharacterHandler(sessionStore, logger, appServices.ConfigService, appServices.DataStore)
	skillPlanHandler := flyHandlers.NewSkillPlanHandler(logger, appServices.SkillService, appServices.DataStore)

	// Define routes
	r.HandleFunc("/callback/", authHandler.CallBack())
	r.HandleFunc("/api/add-character", authHandler.AddCharacterHandler())
	r.HandleFunc("/api/home-data", dashboardHandler.GetDashboardData()).Methods("GET")
	r.HandleFunc("/api/logout", authHandler.Logout())
	r.HandleFunc("/api/login", authHandler.Login())
	r.HandleFunc("/api/get-skill-plan", skillPlanHandler.GetSkillPlanFile())
	r.HandleFunc("/api/save-skill-plan", skillPlanHandler.SaveSkillPlan())
	r.HandleFunc("/api/delete-skill-plan", skillPlanHandler.DeleteSkillPlan())
	r.HandleFunc("/api/unassigned-characters", characterHandler.GetUnassignedCharacters())
	r.HandleFunc("/api/assign-character", characterHandler.AssignCharacter())
	r.HandleFunc("/api/update-account-name", accountHandler.UpdateAccountName())
	r.HandleFunc("/api/update-character", characterHandler.UpdateCharacter())
	r.HandleFunc("/api/toggle-account-status", accountHandler.ToggleAccountStatus())
	r.HandleFunc("/api/remove-character", characterHandler.RemoveCharacter())

	r.HandleFunc("/api/choose-settings-dir", characterHandler.ChooseSettingsDir)
	r.HandleFunc("/api/sync-sub-dir", characterHandler.SyncSubDirectory)
	r.HandleFunc("/api/backup-directory", characterHandler.BackupDirectory)
	r.HandleFunc("/api/associate-character", characterHandler.AssociateCharacter)
	r.HandleFunc("/api/unassociate-character", characterHandler.UnassociateCharacter)
	r.HandleFunc("/api/save-user-selections", characterHandler.SaveUserSelections)
	r.HandleFunc("/api/sync-all-dirs", characterHandler.SyncAllSubdirectories)

	r.HandleFunc("/api/reset-identities", authHandler.ResetAccounts())

	// Serve static files
	staticFileServer := http.FileServer(http.FS(embed.StaticFilesSub))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticFileServer))

	return r
}
