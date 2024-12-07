package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"

	http2 "github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services"
)

func sameIdentities(users []int64, characterIDs []int64) bool {
	if len(characterIDs) != len(users) {
		return false
	}

	for _, userID := range users {
		if !slices.Contains(characterIDs, userID) {
			return false
		}
	}

	return true
}

func sameUserCount(session *sessions.Session, previousUsers, storeUsers int) bool {
	if previousUsers == 0 {
		return false
	}

	if previousUsers != storeUsers {
		return false
	}

	if authenticatedUsers, ok := session.Values[http2.AllAuthenticatedCharacters].([]int64); ok {
		return previousUsers == len(authenticatedUsers)
	}

	return false
}

func handleErrorWithRedirect(w http.ResponseWriter, r *http.Request, redirectURL string) {
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func clearSession(s *http2.SessionService, w http.ResponseWriter, r *http.Request, logger *logrus.Logger) {
	session, err := s.Get(r, http2.SessionName)
	if err != nil {
		logger.Errorf("Failed to get session to clear: %v", err)
	}

	session.Values = make(map[interface{}]interface{})

	err = sessions.Save(r, w)
	if err != nil {
		logger.Errorf("Failed to save session to clear: %v", err)
	}
}

func checkIfCanSkip(session *sessions.Session, sessionValues http2.SessionValues, r *http.Request) (model.AppState, string, bool) {
	canSkip := true
	storeData, etag, ok := persist.Store.Get(sessionValues.LoggedInUser)
	if !ok {
		canSkip = false
	}

	var identitiesCount int
	for _, account := range storeData.Accounts {
		identitiesCount += len(account.Characters)
	}

	if !sameUserCount(session, sessionValues.PreviousUserCount, identitiesCount) {
		canSkip = false
	}
	return storeData, etag, canSkip
}

func getConfigData(logger *logrus.Logger, dataStore *persist.DataStore) model.ConfigData {
	config, err := dataStore.FetchConfigData()
	if err != nil {
		logger.Errorf("error retrieving config data")
		return model.ConfigData{}
	}
	return *config
}

func prepareAppData(
	accounts []model.Account,
	logger *logrus.Logger,
	skill *services.SkillService,
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

func updateStoreAndSession(data model.AppState, etag string, session *sessions.Session, r *http.Request, w http.ResponseWriter) (string, error) {
	newEtag, err := persist.GenerateETag(data)
	if err != nil {
		return etag, fmt.Errorf("failed to generate etag: %w", err)
	}

	loggedIn, ok := session.Values[http2.LoggedInUser].(int64)
	if !ok || loggedIn == 0 {
		return etag, errors.New("user not logged in")
	}

	if newEtag != etag {
		etag, err = persist.Store.Set(loggedIn, data)
		if err != nil {
			return etag, fmt.Errorf("failed to update persist: %w", err)
		}
	}

	session.Values[http2.PreviousEtagUsed] = etag
	if authenticatedUsers, ok := session.Values[http2.AllAuthenticatedCharacters].([]int64); ok {
		session.Values[http2.PreviousUserCount] = len(authenticatedUsers)
	}

	if err = session.Save(r, w); err != nil {
		return etag, fmt.Errorf("failed to save session: %w", err)
	}

	return etag, nil
}
