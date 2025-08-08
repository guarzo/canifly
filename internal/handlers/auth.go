// handlers/auth_handler.go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	flyHttp "github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

type AuthHandler struct {
	sessionService   interfaces.SessionService
	esiAPIService    interfaces.ESIAPIService
	logger           interfaces.Logger
	accountService   interfaces.AccountManagementService
	stateService     interfaces.ConfigurationService
	loginService     interfaces.LoginService
	authClient       interfaces.AuthClient
	cache            interfaces.HTTPCacheService
	wsHub            *WebSocketHub
	characterService interfaces.CharacterService
}

func NewAuthHandler(
	s interfaces.SessionService,
	e interfaces.ESIAPIService,
	l interfaces.Logger,
	accountSvc interfaces.AccountManagementService,
	stateSvc interfaces.ConfigurationService,
	login interfaces.LoginService,
	auth interfaces.AuthClient,
	cache interfaces.HTTPCacheService,
	wsHub *WebSocketHub,
	characterSvc interfaces.CharacterService,
) *AuthHandler {
	return &AuthHandler{
		sessionService:   s,
		esiAPIService:    e,
		logger:           l,
		accountService:   accountSvc,
		stateService:     stateSvc,
		loginService:     login,
		authClient:       auth,
		cache:            cache,
		wsHub:            wsHub,
		characterService: characterSvc,
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

		state, err := h.loginService.GenerateAndStoreInitialState(request.Account)
		if err != nil {
			respondError(w, "Unable to generate state", http.StatusInternalServerError)
			return
		}

		h.logger.Infof("Getting auth URL for state: %s", state)
		url := h.authClient.GetAuthURL(state)
		h.logger.Infof("Generated auth URL: %s", url)
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

		state, err := h.loginService.GenerateAndStoreInitialState(request.Account)
		if err != nil {
			respondError(w, "Unable to generate state", http.StatusInternalServerError)
			return
		}

		url := h.authClient.GetAuthURL(state)
		respondJSON(w, map[string]string{"redirectURL": url, "state": state})
	}
}

func (h *AuthHandler) CallBack() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("callback request received")
		devMode := os.Getenv("DEV_MODE") == "true"

		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		accountName, _, ok := h.loginService.ResolveAccountAndStatusByState(state)
		if !ok {
			h.logger.Error("unable to retrieve value from state")
			handleErrorWithRedirect(w, r, "/")
			return
		}

		h.logger.Infof("Received accountName (account name): %v", accountName)

		token, err := h.authClient.ExchangeCode(code)
		if err != nil {
			h.logger.Errorf("Failed to exchange token for code: %s, %v", code, err)
			handleErrorWithRedirect(w, r, "/")
			return
		}

		user, err := h.esiAPIService.GetUserInfo(token)
		if err != nil {
			h.logger.Errorf("Failed to get user info: %v", err)
			handleErrorWithRedirect(w, r, "/")
			return
		}
		h.logger.Warnf("character is %s", user.CharacterName)

		// Login state is now tracked only via session cookie

		// Use AccountManagementService to handle account creation
		if err = h.accountService.FindOrCreateAccount(accountName, user, token); err != nil {
			h.logger.Errorf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		// Refresh the character to fetch ESI data
		h.logger.Infof("Fetching ESI data for character %s (ID: %d)", user.CharacterName, user.CharacterID)
		if h.characterService != nil {
			updated, err := h.characterService.RefreshCharacterData(user.CharacterID)
			if err != nil {
				h.logger.Errorf("Failed to fetch ESI data: %v", err)
			} else if updated {
				h.logger.Infof("ESI data fetched successfully for %s", user.CharacterName)
			}
		}

		// Clear account cache to ensure fresh data on next fetch
		if h.cache != nil {
			h.cache.Invalidate("accounts:")
			h.logger.Info("Cleared account cache after adding character")
		}

		// Broadcast update via WebSocket
		if h.wsHub != nil {
			h.wsHub.BroadcastUpdate("account:updated", map[string]interface{}{
				"message": "Character added successfully",
			})
			h.logger.Info("Broadcast account update via WebSocket")
		}

		// Mark the state as callback complete (but don't set session here - that happens in finalize-login)
		err = h.loginService.UpdateStateStatusAfterCallBack(state)
		if err != nil {
			respondError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if devMode {
			// In dev, set session cookie and redirect back to dev server
			session, _ := h.sessionService.Get(r, flyHttp.SessionName)
			session.Values[flyHttp.LoggedIn] = true
			if err = session.Save(r, w); err != nil {
				h.logger.Errorf("Error saving session: %v", err)
			}

			frontendPort := os.Getenv("FRONTEND_PORT")
			if frontendPort == "" {
				frontendPort = "3113" // Default to 3113 if not set
			}
			http.Redirect(w, r, "http://localhost:"+frontendPort, http.StatusFound)
			return
		}

		// In production, redirect to static success page
		// Session will be set when Electron app calls finalize-login
		http.Redirect(w, r, "/static/success.html", http.StatusFound)

	}
}

func (h *AuthHandler) FinalizeLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := r.URL.Query().Get("state")
		h.logger.Infof("FinalizeLogin called with state: %s from %s", state, r.RemoteAddr)

		// Debug: List all available states
		h.logger.Debugf("FinalizeLogin: checking state store for state: %s", state)

		accountName, status, ok := h.loginService.ResolveAccountAndStatusByState(state)
		if !ok {
			h.logger.Errorf("FinalizeLogin: state not found: %s", state)
			// Return more informative error for debugging
			respondJSON(w, map[string]interface{}{
				"success": false,
				"error":   "state_not_found",
				"message": fmt.Sprintf("State %s not found in store", state),
			})
			return
		}

		h.logger.Infof("FinalizeLogin: state found - account: %s, callback complete: %v", accountName, status)

		if !status {
			h.logger.Info("FinalizeLogin: callback not yet completed, returning pending status")
			respondJSON(w, map[string]interface{}{
				"success": false,
				"pending": true,
				"message": "Waiting for OAuth callback to complete",
			})
			return
		}

		session, _ := h.sessionService.Get(r, flyHttp.SessionName)
		session.Values[flyHttp.LoggedIn] = true
		h.logger.Infof("FinalizeLogin: Setting session cookie - Domain: %s, Path: /", r.Host)
		if err := session.Save(r, w); err != nil {
			h.logger.Errorf("FinalizeLogin: failed to save session: %v", err)
			respondError(w, "failed to set session", http.StatusInternalServerError)
			return
		}

		// Debug: Check if cookie was actually set
		for _, cookie := range w.Header()["Set-Cookie"] {
			h.logger.Infof("FinalizeLogin: Set-Cookie header: %s", cookie)
		}

		h.logger.Info("FinalizeLogin: session set successfully, login complete")
		//h.loginService.ClearState(state)

		// Return a session token that the frontend can store
		// This is a workaround for file:// protocol not supporting cookies
		sessionToken := state // Use state as a simple token for now

		respondJSON(w, map[string]interface{}{
			"success": true,
			"message": "Login completed successfully",
			"token":   sessionToken,
		})
	}
}

func (h *AuthHandler) GetSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for token in Authorization header (for file:// protocol support)
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			h.logger.Infof("Session check with token: %s", token)

			// Check if this token corresponds to a completed login
			accountName, callbackComplete, ok := h.loginService.ResolveAccountAndStatusByState(token)
			if ok && callbackComplete {
				h.logger.Infof("Token valid - account: %s", accountName)
				respondJSON(w, map[string]interface{}{
					"status":        "ok",
					"authenticated": true,
					"user":          accountName,
				})
				return
			}
			h.logger.Warnf("Invalid or incomplete token: %s", token)
		}

		// Fall back to cookie-based session (for development)
		session, err := h.sessionService.Get(r, flyHttp.SessionName)
		if err != nil {
			// Check if this is a securecookie error (invalid value)
			if strings.Contains(err.Error(), "securecookie: the value is not valid") {
				h.logger.Warnf("Invalid session cookie detected: %v", err)
				// Return unauthenticated status instead of error
				respondJSON(w, map[string]interface{}{
					"status":        "ok",
					"authenticated": false,
				})
				return
			}
			h.logger.Errorf("Failed to get session: %v", err)
			respondError(w, "Failed to get session", http.StatusInternalServerError)
			return
		}

		loggedIn, ok := session.Values[flyHttp.LoggedIn].(bool)
		if !ok {
			loggedIn = false
		}

		h.logger.Infof("Session check - Cookie LoggedIn: %v", loggedIn)

		respondJSON(w, map[string]interface{}{
			"status":        "ok",
			"authenticated": loggedIn,
		})
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
		// Session cleared, no need to clear app state

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
