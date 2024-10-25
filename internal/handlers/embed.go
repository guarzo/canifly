package handlers

import (
	"net/http"
	"path"
	"strings"

	"github.com/gambtho/canifly/internal/embed"
)

// StaticFileHandler serves static files embedded in the binary.
func StaticFileHandler(w http.ResponseWriter, r *http.Request) {
	// Trim the URL path to access the correct file in the static directory
	filePath := strings.TrimPrefix(r.URL.Path, "/static/")
	filePath = path.Join("static", filePath)

	// Serve the embedded file directly
	content, err := embed.StaticFiles.ReadFile(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Set appropriate content type based on file extension (optional)
	w.Header().Set("Content-Type", getContentType(filePath))
	w.Write(content)
}

// Helper function to determine content type based on file extension
func getContentType(filePath string) string {
	ext := path.Ext(filePath)
	switch ext {
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	default:
		return "application/octet-stream"
	}
}
