// http/middleware.go
package http

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

func AuthMiddleware(s interfaces.SessionService, logger interfaces.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(r.URL.Path)
			// Define public routes that don't require authentication
			publicRoutes := map[string]bool{
				"/static":             true,
				"/landing":            true,
				"/api/login":          true,
				"/api/logout":         true,
				"/callback":           true,
				"/api/add-character":  true,
				"/api/finalize-login": true,
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

			logger.WithField("path", r.URL.Path).Debug("Authentication required for private route")

			// Retrieve the session
			session, err := s.Get(r, SessionName)
			if err != nil {
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
