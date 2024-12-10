// handlers/dashboard_handler.go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	flyHttp "github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
	"github.com/guarzo/canifly/internal/services/settings"
)

type DashboardHandler struct {
	sessionService   *flyHttp.SessionService
	skillService     interfaces.SkillService
	characterService interfaces.CharacterService
	accountService   interfaces.AccountService
	settingsService  settings.SettingsService
	stateService     interfaces.StateService
	dataStore        interfaces.DataStore // to retrieve skill plans, skill types, user selections, config
	logger           interfaces.Logger
}

func NewDashboardHandler(
	s *flyHttp.SessionService,
	logger interfaces.Logger,
	skillService interfaces.SkillService,
	characterService interfaces.CharacterService,
	accountService interfaces.AccountService,
	settingsService settings.SettingsService,
	stateService interfaces.StateService,
	dataStore interfaces.DataStore,
) *DashboardHandler {
	return &DashboardHandler{
		sessionService:   s,
		logger:           logger,
		skillService:     skillService,
		characterService: characterService,
		accountService:   accountService,
		settingsService:  settingsService,
		stateService:     stateService,
		dataStore:        dataStore,
	}
}

func (h *DashboardHandler) handleAppStateRefresh(w http.ResponseWriter, noCache bool) {
	storeData := h.stateService.GetAppState()

	if !noCache && len(storeData.Accounts) > 0 {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(storeData)
		go h.refreshDataInBackground()
		return
	}

	updatedData, err := h.refreshAccountsAndState()
	if err != nil {
		h.logger.Errorf("%v", err)
		http.Error(w, `{"error":"Failed to validate accounts"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedData); err != nil {
		http.Error(w, `{"error":"Failed to encode data"}`, http.StatusInternalServerError)
		return
	}
}

func (h *DashboardHandler) refreshAccountsAndState() (model.AppState, error) {
	accounts, err := h.refreshAccounts()
	if err != nil {
		return model.AppState{}, fmt.Errorf("failed to validate accounts: %v", err)
	}

	updatedData := h.prepareAppData(accounts)
	if err = h.updateAndSaveAppState(updatedData); err != nil {
		// Log the error but still try to return updatedData if we have it
		h.logger.Errorf("Failed to update persist and session: %v", err)
	}
	return updatedData, nil
}

func (h *DashboardHandler) updateAndSaveAppState(data model.AppState) error {
	if err := h.stateService.SetAppState(data); err != nil {
		return fmt.Errorf("failed to set app state: %w", err)
	}
	if err := h.stateService.SaveAppStateSnapshot(data); err != nil {
		return fmt.Errorf("failed to save app state snapshot: %w", err)
	}
	return nil
}

func (h *DashboardHandler) GetDashboardData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Infof("GetAppData Called")
		h.handleAppStateRefresh(w, false)
	}
}

func (h *DashboardHandler) GetDashboardDataNoCache() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Infof("GetAppDataNoCache Called")
		h.handleAppStateRefresh(w, true)
	}
}

func (h *DashboardHandler) refreshDataInBackground() {
	start := time.Now()
	h.logger.Debugf("Refreshing data in background...")

	_, err := h.refreshAccountsAndState()
	if err != nil {
		h.logger.Errorf("Failed in background refresh: %v", err)
		return
	}

	timeElapsed := time.Since(start)
	h.logger.Infof("Background refresh complete in %s", timeElapsed)
}

func (h *DashboardHandler) refreshAccounts() ([]model.Account, error) {
	h.logger.Debugf("Refreshing accounts")
	accounts, err := h.accountService.FetchAccounts()
	if err != nil {
		return nil, fmt.Errorf("failed to load accounts: %w", err)
	}

	h.logger.Debugf("Fetched %d accounts", len(accounts))

	for i := range accounts {
		account := &accounts[i]
		h.logger.Debugf("Processing account: %s", account.Name)

		for j := range account.Characters {
			charIdentity := &account.Characters[j]
			h.logger.Debugf("Processing character: %s (ID: %d)", charIdentity.Character.CharacterName, charIdentity.Character.CharacterID)

			updatedCharIdentity, err := h.characterService.ProcessIdentity(charIdentity)
			if err != nil {
				h.logger.Errorf("Failed to process identity for character %d: %v", charIdentity.Character.CharacterID, err)
				continue
			}

			account.Characters[j] = *updatedCharIdentity
		}

		h.logger.Debugf("Account %s has %d characters after processing", account.Name, len(account.Characters))
	}

	// Save the updated accounts
	if err := h.accountService.SaveAccounts(accounts); err != nil {
		return nil, fmt.Errorf("failed to save accounts: %w", err)
	}

	// If DataStore or StateService handles cache saving, call that here if needed
	// Assuming DataStore directly since we haven't abstracted cache operations:
	if c, ok := h.dataStore.(interfaces.CacheRepository); ok {
		_ = c.SaveApiCache()
	}

	return accounts, nil
}

func (h *DashboardHandler) prepareAppData(accounts []model.Account) model.AppState {
	// Retrieve skill plans and skill types from the dataStore through interfaces.DataStore
	skillPlans := h.skillService.GetMatchingSkillPlans(
		accounts,
		h.dataStore.GetSkillPlans(),
		h.dataStore.GetSkillTypes(),
	)

	// Fetch config data
	configData, err := h.dataStore.FetchConfigData()
	if err != nil {
		h.logger.Errorf("Failed to fetch config data: %v", err)
		configData = &model.ConfigData{}
	}

	subDirData, err := h.settingsService.LoadCharacterSettings()
	if err != nil {
		h.logger.Errorf("Failed to load character settings: %v", err)
	}

	userSelections, err := h.dataStore.LoadUserSelections()
	if err != nil {
		h.logger.Warnf("Failed to load user selections, defaulting to empty: %v", err)
		userSelections = make(map[string]model.UserSelection)
	}

	return model.AppState{
		LoggedIn:       true,
		Accounts:       accounts,
		SkillPlans:     skillPlans,
		ConfigData:     *configData,
		SubDirs:        subDirData,
		UserSelections: userSelections,
	}
}
