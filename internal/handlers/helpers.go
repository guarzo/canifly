package handlers

import (
	"github.com/gorilla/sessions"
	"github.com/guarzo/canifly/internal/services/interfaces"
	"github.com/sirupsen/logrus"
	"net/http"

	flyHttp "github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services"
)

func handleErrorWithRedirect(w http.ResponseWriter, r *http.Request, redirectURL string) {
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func clearSession(s *flyHttp.SessionService, w http.ResponseWriter, r *http.Request, logger *logrus.Logger) {
	session, err := s.Get(r, flyHttp.SessionName)
	if err != nil {
		logger.Errorf("Failed to get session to clear: %v", err)
	}

	session.Values = make(map[interface{}]interface{})

	err = sessions.Save(r, w)
	if err != nil {
		logger.Errorf("Failed to save session to clear: %v", err)
	}
}

func prepareAppData(
	accounts []model.Account,
	logger *logrus.Logger,
	skill *interfaces.SkillService,
	dataStore *persist.DataStore,
	configService *services.ConfigService,
) model.AppState {
	// Retrieve skill plans and skill types from the data store instead of global variables
	skillPlans := skill.GetMatchingSkillPlans(
		accounts,
		dataStore.GetSkillPlans(),
		dataStore.GetSkillTypes(),
	)

	// Fetch config data directly from the dataStore
	config, err := dataStore.FetchConfigData()
	if err != nil {
		logger.WithError(err).Error("Failed to fetch config data")
		config = &model.ConfigData{} // fallback to an empty config if needed
	}

	subDirData, err := configService.LoadCharacterSettings()
	if err != nil {
		logger.Errorf("Failed to load character settings: %v", err)
	}
	userSelections, err := dataStore.LoadUserSelections()
	if err != nil {
		logger.Warnf("Failed to load user selections, defaulting to empty: %v", err)
		userSelections = make(map[string]model.UserSelection)
	}

	return model.AppState{
		LoggedIn:       true,
		Accounts:       accounts,
		SkillPlans:     skillPlans,
		ConfigData:     *config,
		SubDirs:        subDirData,
		UserSelections: userSelections,
	}
}
