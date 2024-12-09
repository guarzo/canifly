package handlers

import (
	"encoding/json"
	"github.com/guarzo/canifly/internal/services"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/guarzo/canifly/internal/auth"
	flyHttp "github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services/esi"
)

// AuthHandler holds dependencies needed by auth-related handlers.
type AuthHandler struct {
	sessionService *flyHttp.SessionService
	esiService     esi.ESIService
	logger         *logrus.Logger
	dataStore      *persist.DataStore
	configService  *services.ConfigService
}

func NewAuthHandler(s *flyHttp.SessionService, e esi.ESIService, l *logrus.Logger, data *persist.DataStore, c *services.ConfigService) *AuthHandler {
	return &AuthHandler{
		sessionService: s,
		esiService:     e,
		logger:         l,
		dataStore:      data,
		configService:  c,
	}
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Header().Set("Surrogate-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")

		var request struct {
			Account string `json:"account"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			h.logger.Errorf("Invalid request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if request.Account == "" {
			http.Error(w, "Invalid request body, account must be provided", http.StatusBadRequest)
		}

		url := auth.GetAuthURL(request.Account)

		json.NewEncoder(w).Encode(map[string]string{"redirectURL": url})
	}
}

func (h *AuthHandler) AddCharacterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var request struct {
			Account string `json:"account"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			h.logger.Errorf("Invalid request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		state := "Placeholder"
		if request.Account != "" {
			h.logger.Infof("setting state to be: %v", request.Account)
			state = request.Account
		}

		url := auth.GetAuthURL(state)
		json.NewEncoder(w).Encode(map[string]string{"redirectURL": url})
	}
}

func (h *AuthHandler) CallBack() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("callback request received")

		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state") // This is the account name provided by the user

		h.logger.Infof("Received state (account name): %v", state)

		token, err := auth.ExchangeCode(code)
		if err != nil {
			h.logger.Errorf("Failed to exchange token for code: %s, %v", code, err)
			handleErrorWithRedirect(w, r, "/")
			return
		}

		user, err := h.esiService.GetUserInfo(token)
		if err != nil {
			h.logger.Errorf("Failed to get user info: %v", err)
			handleErrorWithRedirect(w, r, "/")
			return
		}
		h.logger.Warnf("character is %s", user.CharacterName)

		err = h.dataStore.SetAppStateLogin(true)
		if err != nil {
			h.logger.Errorf("Failed to get set appstate: %v", err)
			handleErrorWithRedirect(w, r, "/")
			return
		}

		session, _ := h.sessionService.Get(r, flyHttp.SessionName)
		session.Values[flyHttp.LoggedIn] = true

		err = h.configService.FindOrCreateAccount(state, user, token)
		if err != nil {
			h.logger.Errorf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = session.Save(r, w); err != nil {
			h.logger.Errorf("Error saving session: %v", err)
		}

		http.Redirect(w, r, "http://localhost:5173", http.StatusFound)
	}
}

func (h *AuthHandler) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := h.sessionService.Get(r, flyHttp.SessionName)
		if err != nil {
			http.Error(w, "Failed to get session", http.StatusInternalServerError)
			return
		}

		clearSession(h.sessionService, w, r, h.logger)
		if err := session.Save(r, w); err != nil {
			http.Error(w, "Failed to save session", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}

func (h *AuthHandler) ResetAccounts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.dataStore.DeleteAccounts()
		if err != nil {
			h.logger.Errorf("Failed to delete identity %d", err)
		}

		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}
}
