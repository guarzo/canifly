// handlers/utils.go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"

	flyHttp "github.com/guarzo/canifly/internal/http"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

// respondJSON sends a success response with JSON-encoded data.
func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}

// respondError sends an error response in JSON format.
func respondError(w http.ResponseWriter, msg string, code int) {
	http.Error(w, fmt.Sprintf(`{"error":"%s"}`, msg), code)
}

// decodeJSONBody decodes the JSON body into the provided dst.
func decodeJSONBody(r *http.Request, dst interface{}) error {
	return json.NewDecoder(r.Body).Decode(dst)
}

func clearSession(s *flyHttp.SessionService, w http.ResponseWriter, r *http.Request, logger interfaces.Logger) {
	session, err := s.Get(r, flyHttp.SessionName)
	if err != nil {
		logger.Errorf("Failed to get session to clear: %v", err)
	}

	session.Values = make(map[interface{}]interface{})

	err = sessions.Save(r, w)
	if err != nil {
		logger.Errorf("Failed to save session to clear: %v", err)
	}
}
func handleErrorWithRedirect(w http.ResponseWriter, r *http.Request, redirectURL string) {
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
