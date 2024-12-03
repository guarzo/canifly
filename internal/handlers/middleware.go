package handlers

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/guarzo/canifly/internal/utils/xlog"
)

func AuthMiddleware(s *SessionService) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// List of public routes that don't require authentication
			publicRoutes := map[string]bool{
				"/static":   true,
				"/landing":  true,
				"/login":    true,
				"/logout":   true,
				"/callback": true,
			}

			// Check if the path starts with one of the public routes
			for publicRoute := range publicRoutes {
				if strings.HasPrefix(r.URL.Path, publicRoute) {
					next.ServeHTTP(w, r)
					return
				}
			}

			xlog.Log("Proceeding to authentication check")

			session, _ := s.Get(r, sessionName)
			sessionValues := getSessionValues(session)

			// Check if logged_in_user is present
			currentUser := sessionValues.LoggedInUser
			if currentUser == 0 {
				xlog.Logf("User not logged in, redirecting to /landing")
				http.Error(w, `{"error":"user is not logged in"}`, http.StatusUnauthorized)
			}

			next.ServeHTTP(w, r)
		})
	}
}
