// handlers/auth_handler.go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	flyHttp "github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/persist"
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
	sessionStore     *persist.SessionStore // Add persistent session store
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
	sessionStore *persist.SessionStore,
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
		sessionStore:     sessionStore,
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

		// IMPORTANT: Clear cache and broadcast AFTER ESI data is fetched and saved
		// This ensures the frontend gets the complete data including ESI information

		// Clear account cache to ensure fresh data on next fetch
		if h.cache != nil {
			// Invalidate all account-related cache entries with various pagination params
			h.cache.Invalidate("accounts:")
			// Also specifically invalidate common cache keys
			h.cache.Invalidate("accounts:list:page:1:limit:20")
			h.cache.Invalidate("accounts:list:page:1:limit:1000")
			h.logger.Infof("Invalidated all account cache entries after adding character %s", user.CharacterName)
		}

		// Broadcast update via WebSocket - frontend will fetch fresh data
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

		// Create persistent session if session store is available
		var sessionToken string
		if h.sessionStore != nil {
			// Get character IDs for this account
			var characterIDs []int64
			accounts, err := h.accountService.FetchAccounts()
			if err == nil {
				for _, acc := range accounts {
					if acc.Name == accountName {
						for _, char := range acc.Characters {
							characterIDs = append(characterIDs, char.Character.CharacterID)
						}
						break
					}
				}
			}

			// Check if client wants to remember the session
			rememberMe := r.URL.Query().Get("remember") == "true"

			// Create persistent session
			sessionData, err := h.sessionStore.CreateSession(accountName, characterIDs, rememberMe)
			if err != nil {
				h.logger.Errorf("Failed to create persistent session: %v", err)
				// Fall back to using state as token
				sessionToken = state
			} else {
				sessionToken = sessionData.SessionID
				h.logger.Infof("Created persistent session %s for account %s", sessionToken, accountName)
			}
		} else {
			// Fall back to using state as token
			sessionToken = state
		}

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

			// First check persistent session store
			if h.sessionStore != nil {
				if sessionData, exists := h.sessionStore.GetSession(token); exists {
					h.logger.Infof("Valid persistent session found for account: %s", sessionData.AccountName)
					respondJSON(w, map[string]interface{}{
						"status":        "ok",
						"authenticated": true,
						"user":          sessionData.AccountName,
						"sessionId":     sessionData.SessionID,
					})
					return
				}
			}

			// Fall back to checking OAuth state (for backward compatibility)
			accountName, callbackComplete, ok := h.loginService.ResolveAccountAndStatusByState(token)
			if ok && callbackComplete {
				h.logger.Infof("Token valid (OAuth state) - account: %s", accountName)
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

// ValidateSession checks if a session is valid and returns detailed info
func (h *AuthHandler) ValidateSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			respondJSON(w, map[string]interface{}{
				"valid":   false,
				"message": "No token provided",
			})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Check persistent session store
		if h.sessionStore != nil {
			if sessionData, exists := h.sessionStore.GetSession(token); exists {
				// Validate that at least one account exists (basic sanity check)
				accounts, err := h.accountService.FetchAccounts()
				if err == nil && len(accounts) > 0 {
					// Check if the account still exists
					accountExists := false
					for _, acc := range accounts {
						if acc.Name == sessionData.AccountName {
							accountExists = true
							break
						}
					}

					if !accountExists {
						h.logger.Warnf("Session references non-existent account: %s", sessionData.AccountName)
						h.sessionStore.DeleteSession(token)
						respondJSON(w, map[string]interface{}{
							"valid":   false,
							"message": "Account not found",
						})
						return
					}
				}

				// Return session info
				respondJSON(w, map[string]interface{}{
					"valid":       true,
					"sessionId":   sessionData.SessionID,
					"accountName": sessionData.AccountName,
					"characters":  sessionData.Characters,
					"createdAt":   sessionData.CreatedAt,
					"expiresAt":   sessionData.ExpiresAt,
					"rememberMe":  sessionData.RememberMe,
				})
				return
			}
		}

		// Fall back to OAuth state check
		accountName, callbackComplete, ok := h.loginService.ResolveAccountAndStatusByState(token)
		if ok && callbackComplete {
			respondJSON(w, map[string]interface{}{
				"valid":       true,
				"accountName": accountName,
				"legacy":      true, // Indicate this is a legacy session
			})
			return
		}

		respondJSON(w, map[string]interface{}{
			"valid":   false,
			"message": "Invalid or expired session",
		})
	}
}

// RefreshSession generates a new session ID while preserving session data
func (h *AuthHandler) RefreshSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get current token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			respondError(w, "No token provided", http.StatusUnauthorized)
			return
		}

		oldToken := strings.TrimPrefix(authHeader, "Bearer ")

		if h.sessionStore == nil {
			respondError(w, "Session refresh not available", http.StatusServiceUnavailable)
			return
		}

		// Refresh the session
		newSession, err := h.sessionStore.RefreshSession(oldToken)
		if err != nil {
			h.logger.Errorf("Failed to refresh session: %v", err)
			respondError(w, "Failed to refresh session", http.StatusInternalServerError)
			return
		}

		respondJSON(w, map[string]interface{}{
			"success":   true,
			"token":     newSession.SessionID,
			"expiresAt": newSession.ExpiresAt,
		})
	}
}

// GetActiveSessions returns all active sessions for the current account
func (h *AuthHandler) GetActiveSessions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get current session to identify account
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			respondError(w, "Not authenticated", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		if h.sessionStore == nil {
			respondJSON(w, map[string]interface{}{
				"sessions": []interface{}{},
			})
			return
		}

		// Get current session to find account name
		currentSession, exists := h.sessionStore.GetSession(token)
		if !exists {
			respondError(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		// Get all sessions for this account
		sessions := h.sessionStore.GetAccountSessions(currentSession.AccountName)

		// Format response
		sessionList := make([]map[string]interface{}, len(sessions))
		for i, s := range sessions {
			sessionList[i] = map[string]interface{}{
				"sessionId":    s.SessionID,
				"createdAt":    s.CreatedAt,
				"lastAccessed": s.LastAccessed,
				"expiresAt":    s.ExpiresAt,
				"current":      s.SessionID == token,
			}
		}

		respondJSON(w, map[string]interface{}{
			"sessions": sessionList,
		})
	}
}
