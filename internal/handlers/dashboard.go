package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"

	http2 "github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services"
	"github.com/guarzo/canifly/internal/services/esi"
)

// DashboardHandler holds the dependencies needed by the home-related handlers.
type DashboardHandler struct {
	sessionService *http2.SessionService
	esiService     esi.ESIService
	skillService   *services.SkillService
	dataStore      *persist.DataStore
	logger         *logrus.Logger
	configService  *services.ConfigService
}

// NewDashboardHandler creates a new DashboardHandler with the given session and ESI services.
func NewDashboardHandler(s *http2.SessionService, e esi.ESIService, l *logrus.Logger, skill *services.SkillService, data *persist.DataStore, configService *services.ConfigService) *DashboardHandler {
	return &DashboardHandler{
		sessionService: s,
		esiService:     e,
		logger:         l,
		skillService:   skill,
		dataStore:      data,
		configService:  configService,
	}
}

func (h *DashboardHandler) GetDashboardData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := h.sessionService.Get(r, http2.SessionName)
		sessionValues := http2.GetSessionValues(session)

		if sessionValues.LoggedInUser == 0 {
			http.Error(w, `{"error":"Failed to encode data"}`, http.StatusUnauthorized)
			return
		}

		storeData, etag, canSkip := checkIfCanSkip(session, sessionValues, r)

		if canSkip {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(storeData); err != nil {
				http.Error(w, `{"error":"Failed to encode data"}`, http.StatusInternalServerError)
			}
			return
		}

		accounts, err := h.validateAccounts(session, storeData)
		if err != nil {
			h.logger.Errorf("Failed to validate accounts: %v", err)
			http.Error(w, `{"error":"Failed to validate accounts"}`, http.StatusInternalServerError)
			return
		}

		data := prepareAppData(accounts, h.logger, h.skillService, h.dataStore, h.configService)

		_, err = updateStoreAndSession(data, etag, session, r, w)
		if err != nil {
			h.logger.Errorf("Failed to update persist and session: %v", err)
			http.Error(w, `{"error":"Failed to update session"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, `{"error":"Failed to encode data"}`, http.StatusInternalServerError)
			return
		}
	}
}

// validateAccounts is now a method on DashboardHandler, allowing access to h.esiService.
// It replaces the direct calls to esi.ProcessIdentity with h.esiService.ProcessIdentity.
func (h *DashboardHandler) validateAccounts(session *sessions.Session, storeData model.AppState) ([]model.Account, error) {
	authenticatedUsers, ok := session.Values[http2.AllAuthenticatedCharacters].([]int64)
	if !ok {
		h.logger.Errorf("Failed to retrieve authenticated users from session")
		return nil, fmt.Errorf("failed to retrieve authenticated users from session")
	}

	needIdentityPopulation := len(authenticatedUsers) == 0
	if len(storeData.Accounts) != 0 {
		for _, account := range storeData.Accounts {
			var accountCharacterIDs []int64
			for _, charIdentity := range account.Characters {
				accountCharacterIDs = append(accountCharacterIDs, charIdentity.Character.CharacterID)
			}
			needIdentityPopulation = needIdentityPopulation || !sameIdentities(authenticatedUsers, accountCharacterIDs)
		}
	} else {
		needIdentityPopulation = true
	}

	if needIdentityPopulation {
		h.logger.Infof("Need to populate identities")
		accounts, err := h.dataStore.FetchAccounts()
		if err != nil {
			h.logger.Errorf("Failed to load accounts: %v", err)
			return nil, fmt.Errorf("failed to load accounts: %w", err)
		}

		h.logger.Infof("Fetched %d accounts", len(accounts))

		// Process each account and its characters using the injected ESIService
		for i := range accounts {
			account := &accounts[i]
			h.logger.Infof("Processing account: %s", account.Name)

			for j := range account.Characters {
				charIdentity := &account.Characters[j]
				h.logger.Infof("Processing character: %s (ID: %d)", charIdentity.Character.CharacterName, charIdentity.Character.CharacterID)

				// Use h.esiService instead of directly calling esi.ProcessIdentity
				updatedCharIdentity, err := h.esiService.ProcessIdentity(charIdentity)
				if err != nil {
					h.logger.Errorf("Failed to process identity for character %d: %v", charIdentity.Character.CharacterID, err)
					continue
				}

				account.Characters[j] = *updatedCharIdentity
			}

			h.logger.Infof("Account %s has %d characters after processing", account.Name, len(account.Characters))
		}

		// Save the updated accounts
		if err := h.dataStore.SaveAccounts(accounts); err != nil {
			h.logger.Errorf("Failed to save accounts: %v", err)
			return nil, fmt.Errorf("failed to save accounts: %w", err)
		}

		storeData.Accounts = accounts

		var allCharacterIDs []int64
		for _, account := range accounts {
			for _, charIdentity := range account.Characters {
				allCharacterIDs = append(allCharacterIDs, charIdentity.Character.CharacterID)
			}
		}

		session.Values[http2.AllAuthenticatedCharacters] = allCharacterIDs
		_ = h.dataStore.SaveApiCache()
	}

	return storeData.Accounts, nil
}
