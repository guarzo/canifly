package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	flyHttp "github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services"
	"github.com/guarzo/canifly/internal/services/esi"
)

// DashboardHandler holds the dependencies needed by the home-related handlers.
type DashboardHandler struct {
	sessionService *flyHttp.SessionService
	esiService     esi.ESIService
	skillService   *services.SkillService
	dataStore      *persist.DataStore
	logger         *logrus.Logger
	configService  *services.ConfigService
}

// NewDashboardHandler creates a new DashboardHandler with the given session and ESI services.
func NewDashboardHandler(s *flyHttp.SessionService, e esi.ESIService, l *logrus.Logger, skill *services.SkillService, data *persist.DataStore, configService *services.ConfigService) *DashboardHandler {
	return &DashboardHandler{
		sessionService: s,
		esiService:     e,
		logger:         l,
		skillService:   skill,
		dataStore:      data,
		configService:  configService,
	}
}

func (h *DashboardHandler) handleAppStateRefresh(w http.ResponseWriter, noCache bool) {
	// If noCache is true, we directly refresh from scratch
	// If noCache is false, we can check if we have cached data first

	storeData := h.dataStore.GetAppState()

	// If we are not forcing no-cache and we have cached data
	// immediately return cached data and trigger background refresh
	if !noCache && len(storeData.Accounts) > 0 {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(storeData)
		go h.refreshDataInBackground()
		return
	}

	updatedData, err := h.refreshAccountsAndState(storeData)
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

func (h *DashboardHandler) refreshAccountsAndState(storeData model.AppState) (model.AppState, error) {
	err := h.refreshAccounts(storeData)
	if err != nil {
		return model.AppState{}, fmt.Errorf("failed to validate accounts: %v", err)
	}

	updatedData := prepareAppData(storeData.Accounts, h.logger, h.skillService, h.dataStore, h.configService)
	if err = h.updateAndSaveAppState(updatedData); err != nil {
		// Log the error but still try to return updatedData if we have it
		h.logger.Errorf("Failed to update persist and session: %v", err)
	}
	return updatedData, nil
}

func (h *DashboardHandler) updateAndSaveAppState(data model.AppState) error {
	h.dataStore.SetAppState(data)
	if err := h.dataStore.SaveAppStateSnapshot(data); err != nil {
		return fmt.Errorf("failed to save app state snapshot: %w", err)
	}

	return nil
}

func (h *DashboardHandler) GetDashboardData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.handleAppStateRefresh(w, false)
	}
}

func (h *DashboardHandler) GetDashboardDataNoCache() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.handleAppStateRefresh(w, true)
	}
}

func (h *DashboardHandler) refreshDataInBackground() {
	start := time.Now()
	h.logger.Debugf("Refreshing data in background...")
	storeData := h.dataStore.GetAppState()

	_, err := h.refreshAccountsAndState(storeData)
	if err != nil {
		h.logger.Errorf("Failed in background refresh: %v", err)
		return
	}

	timeElapsed := time.Since(start)
	h.logger.Infof("Background refresh complete in %s", timeElapsed)
}

func (h *DashboardHandler) refreshAccounts(storeData model.AppState) error {
	h.logger.Debugf("Refreshing accounts ")
	accounts, err := h.dataStore.FetchAccounts()
	if err != nil {
		return fmt.Errorf("failed to load accounts: %w", err)
	}

	h.logger.Debugf("Fetched %d accounts", len(accounts))

	// Process each account and its characters using the injected ESIService
	for i := range accounts {
		account := &accounts[i]
		h.logger.Debugf("Processing account: %s", account.Name)

		for j := range account.Characters {
			charIdentity := &account.Characters[j]
			h.logger.Debugf("Processing character: %s (ID: %d)", charIdentity.Character.CharacterName, charIdentity.Character.CharacterID)

			// Use h.esiService instead of directly calling esi.ProcessIdentity
			updatedCharIdentity, err := h.esiService.ProcessIdentity(charIdentity)
			if err != nil {
				h.logger.Errorf("Failed to process identity for character %d: %v", charIdentity.Character.CharacterID, err)
				continue
			}

			account.Characters[j] = *updatedCharIdentity
		}

		h.logger.Debugf("Account %s has %d characters after processing", account.Name, len(account.Characters))
	}

	// Save the updated accounts
	if err := h.dataStore.SaveAccounts(accounts); err != nil {
		return fmt.Errorf("failed to save accounts: %w", err)
	}

	storeData.Accounts = accounts
	_ = h.dataStore.SaveApiCache()

	return nil
}
