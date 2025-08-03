package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/guarzo/canifly/internal/services/interfaces"
	"io"
	"net/http"
)

// HandleServiceError handles service errors with consistent logging and response
func HandleServiceError(w http.ResponseWriter, logger interfaces.Logger, err error, action string) {
	logger.Errorf("Failed to %s: %v", action, err)
	respondError(w, fmt.Sprintf("Failed to %s", action), http.StatusInternalServerError)
}

// DecodeAndValidate decodes JSON request body and returns typed result
func DecodeAndValidate[T any](r *http.Request, w http.ResponseWriter) (*T, bool) {
	var req T
	if err := decodeJSONBody(r, &req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return nil, false
	}
	return &req, true
}

// HandleNotFound sends a standardized not found response
func HandleNotFound(w http.ResponseWriter, logger interfaces.Logger, resource string) {
	logger.Debugf("%s not found", resource)
	respondError(w, fmt.Sprintf("%s not found", resource), http.StatusNotFound)
}

// HandleUnauthorized sends a standardized unauthorized response
func HandleUnauthorized(w http.ResponseWriter, logger interfaces.Logger, reason string) {
	logger.Warnf("Unauthorized access: %s", reason)
	respondError(w, "Unauthorized", http.StatusUnauthorized)
}

// HandleBadRequest sends a standardized bad request response
func HandleBadRequest(w http.ResponseWriter, logger interfaces.Logger, message string) {
	logger.Debugf("Bad request: %s", message)
	respondError(w, message, http.StatusBadRequest)
}

// Basic utility functions used across handlers

// respondError sends a JSON error response
func respondError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// respondJSON sends a JSON success response
func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// decodeJSONBody decodes the JSON request body into the provided interface
func decodeJSONBody(r *http.Request, v interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("request body is nil")
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	if len(body) == 0 {
		return fmt.Errorf("request body is empty")
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	return nil
}

// handleErrorWithRedirect handles errors by redirecting to a given path
func handleErrorWithRedirect(w http.ResponseWriter, r *http.Request, redirectPath string) {
	http.Redirect(w, r, redirectPath, http.StatusSeeOther)
}

// clearSession clears the session data
func clearSession(sessionService interfaces.SessionService, w http.ResponseWriter, r *http.Request, logger interfaces.Logger) {
	session, err := sessionService.Get(r, "canifly-session")
	if err != nil {
		logger.WithError(err).Error("Failed to get session for clearing")
		return
	}

	// Clear session values
	for key := range session.Values {
		delete(session.Values, key)
	}

	// Set MaxAge to -1 to delete the cookie
	session.Options.MaxAge = -1

	if err := session.Save(r, w); err != nil {
		logger.WithError(err).Error("Failed to save cleared session")
	}
}

// respondEncodedData sends encoded data as JSON response
func respondEncodedData(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
