// http/middleware.go
package http

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

// PersistentSessionValidator is the minimal surface AuthMiddleware needs from
// persist.SessionStore. Defining it here avoids an import cycle and lets the
// middleware fall back gracefully when no persistent store is configured.
type PersistentSessionValidator interface {
	ValidateSession(sessionID string) bool
}

func AuthMiddleware(s interfaces.SessionService, loginSvc interfaces.LoginService, persistentSessions PersistentSessionValidator, logger interfaces.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(r.URL.Path)
			// Define public routes that don't require authentication
			publicRoutes := map[string]bool{
				"/health":                     true, // Health check endpoint
				"/static":                     true,
				"/landing":                    true,
				"/api/login":                  true,
				"/api/logout":                 true,
				"/api/session":                true,
				"/callback":                   true,
				"/api/add-character":          true,
				"/api/finalize-login":         true,
				"/api/config/eve/status":      true, // Check if EVE config is needed
				"/api/config/eve/credentials": true, // Save EVE credentials during first-run
				"/api/ws":                     true, // WebSocket endpoint for real-time updates
			}

			// Allow access if the request matches a public route
			for publicRoute := range publicRoutes {
				if strings.HasPrefix(r.URL.Path, publicRoute) {
					logger.WithFields(logrus.Fields{
						"path":   r.URL.Path,
						"public": true,
					}).Debug("Public route accessed")
					next.ServeHTTP(w, r)
					return
				}
			}

			// Check for Authorization header with token (for file:// protocol support)
			authHeader := r.Header.Get("Authorization")

			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				token := strings.TrimPrefix(authHeader, "Bearer ")

				// Prefer the persistent session store: FinalizeLogin issues
				// session IDs from there, so production (file://) clients send
				// those as their bearer token.
				if persistentSessions != nil && persistentSessions.ValidateSession(token) {
					logger.WithFields(logrus.Fields{
						"path":   r.URL.Path,
						"method": "session",
					}).Debug("Session token authentication successful")
					next.ServeHTTP(w, r)
					return
				}

				// Fall back to OAuth login state (legacy / pre-persistent tokens).
				accountName, callbackComplete, ok := loginSvc.ResolveAccountAndStatusByState(token)
				if ok && callbackComplete {
					logger.WithFields(logrus.Fields{
						"path":    r.URL.Path,
						"account": accountName,
						"method":  "token",
					}).Debug("Token authentication successful")
					next.ServeHTTP(w, r)
					return
				}
				logger.WithField("path", r.URL.Path).Warn("Invalid or incomplete token")
				http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			logger.WithField("path", r.URL.Path).Debug("Authentication required for private route")

			// Retrieve the session
			session, err := s.Get(r, SessionName)
			if err != nil {
				// Check if this is a securecookie error (invalid value)
				if strings.Contains(err.Error(), "securecookie: the value is not valid") {
					logger.WithError(err).Warn("Invalid session cookie detected, treating as unauthenticated")
					// Treat as unauthenticated rather than error
					http.Error(w, `{"error":"user is not logged in"}`, http.StatusUnauthorized)
					return
				}
				logger.WithError(err).Error("Failed to retrieve session")
				http.Error(w, `{"error":"failed to retrieve session"}`, http.StatusInternalServerError)
				return
			}

			loggedIn, ok := session.Values[LoggedIn].(bool)
			if !ok || !loggedIn {
				logger.Warn("Unauthenticated access attempt")
				http.Error(w, `{"error":"user is not logged in"}`, http.StatusUnauthorized)
				return
			}

			logger.WithFields(logrus.Fields{
				"path": r.URL.Path,
			}).Debug("User authenticated")

			// Proceed to the next handler
			next.ServeHTTP(w, r)
		})
	}
}
