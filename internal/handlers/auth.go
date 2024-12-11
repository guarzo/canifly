// handlers/auth_handler.go
package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/guarzo/canifly/internal/auth"
	flyHttp "github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

type AuthHandler struct {
	sessionService *flyHttp.SessionService
	esiService     interfaces.ESIService
	logger         interfaces.Logger
	accountService interfaces.AccountService
	stateService   interfaces.StateService
	loginService   interfaces.LoginService
}

func NewAuthHandler(
	s *flyHttp.SessionService,
	e interfaces.ESIService,
	l interfaces.Logger,
	accountSvc interfaces.AccountService,
	stateSvc interfaces.StateService,
	login interfaces.LoginService,
) *AuthHandler {
	return &AuthHandler{
		sessionService: s,
		esiService:     e,
		logger:         l,
		accountService: accountSvc,
		stateService:   stateSvc,
		loginService:   login,
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
			respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if request.Account == "" {
			respondError(w, "account must be provided", http.StatusBadRequest)
			return
		}

		state, err := h.loginService.GenerateAnStoreState(request.Account)
		if err != nil {
			respondError(w, "Unable to generate state", http.StatusInternalServerError)
			return
		}

		url := auth.GetAuthURL(state)
		respondJSON(w, map[string]string{"redirectURL": url, "state": state})
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
			respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if request.Account == "" {
			respondError(w, "Invalid request body - account must be provided", http.StatusBadRequest)
			return
		}

		state, err := h.loginService.GenerateAnStoreState(request.Account)
		if err != nil {
			respondError(w, "Unable to generate state", http.StatusInternalServerError)
			return
		}

		url := auth.GetAuthURL(state)
		respondJSON(w, map[string]string{"redirectURL": url})
	}
}

func (h *AuthHandler) CallBack() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("callback request received")
		devMode := os.Getenv("DEV_MODE") == "true"

		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		accountName, ok := h.loginService.ResolveAccountByState(state)
		if !ok {
			h.logger.Error("unable to retrieve value from state")
			handleErrorWithRedirect(w, r, "/")
			return
		}

		h.logger.Infof("Received accountName (account name): %v", accountName)

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

		if err := h.stateService.SetAppStateLogin(true); err != nil {
			h.logger.Errorf("Failed to set app state: %v", err)
			handleErrorWithRedirect(w, r, "/")
			return
		}

		// Use AccountService to handle account creation
		if err := h.accountService.FindOrCreateAccount(accountName, user, token); err != nil {
			h.logger.Errorf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if devMode {
			// In dev, redirect back to dev server so the internal flow works as before
			session, _ := h.sessionService.Get(r, flyHttp.SessionName)
			session.Values[flyHttp.LoggedIn] = true
			if err = session.Save(r, w); err != nil {
				h.logger.Errorf("Error saving session: %v", err)
			}
			http.Redirect(w, r, "http://localhost:5173", http.StatusFound)
		}

		http.Redirect(w, r, "http://localhost:8713/static/success.html", http.StatusFound)

	}
}

func (h *AuthHandler) FinalizeLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := r.URL.Query().Get("state")
		_, ok := h.loginService.ResolveAccountByState(state)
		if !ok {
			http.Error(w, `{"error":"invalid state or login not completed"}`, http.StatusUnauthorized)
			return
		}

		session, _ := h.sessionService.Get(r, flyHttp.SessionName)
		session.Values[flyHttp.LoggedIn] = true
		if err := session.Save(r, w); err != nil {
			http.Error(w, `{"error":"failed to set session"}`, http.StatusInternalServerError)
			return
		}

		h.loginService.ClearState(state)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true}`))
	}
}

func (h *AuthHandler) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := h.sessionService.Get(r, flyHttp.SessionName)
		if err != nil {
			respondError(w, "Failed to get session", http.StatusInternalServerError)
			return
		}

		clearSession(h.sessionService, w, r, h.logger)
		if err := session.Save(r, w); err != nil {
			respondError(w, "Failed to save session", http.StatusInternalServerError)
			return
		}

		respondJSON(w, map[string]bool{"success": true})
	}
}

func (h *AuthHandler) ResetAccounts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.accountService.DeleteAllAccounts()
		if err != nil {
			h.logger.Errorf("Failed to delete identity %v", err)
		}
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	}
}
