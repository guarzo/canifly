// http/middleware.go
package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func AuthMiddleware(s *SessionService, logger *logrus.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Define public routes that don't require authentication
			publicRoutes := map[string]bool{
				"/static":         true,
				"/landing":        true,
				"/login":          true,
				"/logout":         true,
				"/callback":       true,
				"/auth-character": true,
			}

			// Allow access if the request matches a public route
			for publicRoute := range publicRoutes {
				if strings.HasPrefix(r.URL.Path, publicRoute) {
					logger.WithFields(logrus.Fields{
						"path":   r.URL.Path,
						"public": true,
					}).Info("Public route accessed")
					next.ServeHTTP(w, r)
					return
				}
			}

			logger.WithField("path", r.URL.Path).Info("Authentication required for private route")

			// Retrieve the session
			session, err := s.Get(r, SessionName)
			if err != nil {
				logger.WithError(err).Error("Failed to retrieve session")
				http.Error(w, `{"error":"failed to retrieve session"}`, http.StatusInternalServerError)
				return
			}

			// Get logged-in user from session
			sessionValues := GetSessionValues(session)
			currentUser := sessionValues.LoggedInUser
			if currentUser == 0 {
				logger.Warn("Unauthenticated access attempt")
				http.Error(w, `{"error":"user is not logged in"}`, http.StatusUnauthorized)
				return
			}

			logger.WithFields(logrus.Fields{
				"user": currentUser,
				"path": r.URL.Path,
			}).Info("User authenticated")

			// Add user info to request context
			r = r.WithContext(setUserInContext(r.Context(), currentUser))

			// Proceed to the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// setUserInContext adds the logged-in user to the request context
func setUserInContext(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, LoggedInUser, userID)
}

// getUserFromContext retrieves the logged-in user from the request context
func getUserFromContext(ctx context.Context) int64 {
	if user, ok := ctx.Value(LoggedInUser).(int64); ok {
		return user
	}
	return 0
}
