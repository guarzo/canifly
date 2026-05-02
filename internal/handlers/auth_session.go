package handlers

import (
	"net/http"
	"strings"

	flyHttp "github.com/guarzo/canifly/internal/http"
)

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
