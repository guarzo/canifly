package handlers

import (
	"net/http"
	"os"
	"path"
	"strings"
)

// Import the necessary packages
import (
	"log"
	"mime"
	"path/filepath"
)

func StaticFileHandler(w http.ResponseWriter, r *http.Request) {
	// Log the requested URL path
	log.Printf("Requested URL: %s", r.URL.Path)

	// Extract the path after /static/
	path := strings.TrimPrefix(r.URL.Path, "/static/")
	path = filepath.Clean(path)

	// Build the file path relative to your static files directory
	filePath := filepath.Join("static", path)

	// Log the resolved file path
	log.Printf("Serving file: %s", filePath)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("File not found: %s", filePath)
		http.NotFound(w, r)
		return
	}

	// Determine the Content-Type based on the file extension
	ext := filepath.Ext(filePath)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Log the Content-Type being sent
	log.Printf("Content-Type: %s", contentType)
	w.Header().Set("Content-Type", contentType)

	// Serve the file
	http.ServeFile(w, r, filePath)
}

//
//
//
//// StaticFileHandler serves static files embedded in the binary.
//func StaticFileHandler(w http.ResponseWriter, r *http.Request) {
//	// Trim the URL path to access the correct file in the static directory
//	filePath := strings.TrimPrefix(r.URL.Path, "/static/")
//	filePath = path.Join("static", filePath)
//
//	// Serve the embedded file directly
//	content, err := embed.StaticFiles.ReadFile(filePath)
//	if err != nil {
//		http.NotFound(w, r)
//		return
//	}
//
//	// Set appropriate content type based on file extension (optional)
//	w.Header().Set("Content-Type", getContentType(filePath))
//	w.Write(content)
//}

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
