package server

import (
	"github.com/guarzo/canifly/internal/services/interfaces"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/guarzo/canifly/internal/embed"
	flyHandlers "github.com/guarzo/canifly/internal/handlers"
	flyHttp "github.com/guarzo/canifly/internal/http"
)

// SetupRouter configures and returns the app’s router
func SetupRouter(secret string, logger interfaces.Logger, appServices *AppServices) *mux.Router {
	sessionStore := flyHttp.NewSessionService(secret)
	r := mux.NewRouter()

	// Add authentication middleware
	r.Use(flyHttp.AuthMiddleware(sessionStore, logger))
	dashboardHandler := flyHandlers.NewDashboardHandler(sessionStore, logger, appServices.DashBoardService)
	authHandler := flyHandlers.NewAuthHandler(sessionStore, appServices.EsiService, logger, appServices.AccountService, appServices.StateService, appServices.LoginService)
	accountHandler := flyHandlers.NewAccountHandler(sessionStore, logger, appServices.AccountService)
	characterHandler := flyHandlers.NewCharacterHandler(logger, appServices.CharacterService)
	skillPlanHandler := flyHandlers.NewSkillPlanHandler(logger, appServices.SkillService)
	settingsHandler := flyHandlers.NewSettingsHandler(logger, appServices.SettingsService)
	assocHandler := flyHandlers.NewAssociationHandler(logger, appServices.AssocService)

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

	r.HandleFunc("/api/update-account-name", accountHandler.UpdateAccountName())
	r.HandleFunc("/api/toggle-account-status", accountHandler.ToggleAccountStatus())
	r.HandleFunc("/api/remove-account", accountHandler.RemoveAccount())

	r.HandleFunc("/api/update-character", characterHandler.UpdateCharacter)
	r.HandleFunc("/api/remove-character", characterHandler.RemoveCharacter)

	r.HandleFunc("/api/choose-settings-dir", settingsHandler.ChooseSettingsDir)
	r.HandleFunc("/api/reset-to-default-directory", settingsHandler.ResetToDefaultDir)
	r.HandleFunc("/api/backup-directory", settingsHandler.BackupDirectory)
	r.HandleFunc("/api/save-user-selections", settingsHandler.SaveUserSelections)

	r.HandleFunc("/api/sync-subdirectory", settingsHandler.SyncSubDirectory)
	r.HandleFunc("/api/sync-all-subdirectories", settingsHandler.SyncAllSubdirectories)

	r.HandleFunc("/api/associate-character", assocHandler.AssociateCharacter)
	r.HandleFunc("/api/unassociate-character", assocHandler.UnassociateCharacter)

	// Serve static files
	staticFileServer := http.FileServer(http.FS(embed.StaticFilesSub))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticFileServer))

	return r
}
