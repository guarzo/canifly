// handlers/utils.go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
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
