package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

// WithCache wraps a handler function with caching support
func WithCache(
	cache interfaces.HTTPCacheService,
	logger interfaces.Logger,
	key string,
	ttl time.Duration,
	fetcher func() (interface{}, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check cache first
		if data, etag, found := cache.Get(key); found {
			w.Header().Set("ETag", etag)
			w.Header().Set("Cache-Control", "private, max-age=300") // 5 minutes browser cache
			
			// Check if client has current version
			if r.Header.Get("If-None-Match") == etag {
				w.WriteHeader(http.StatusNotModified)
				return
			}
			
			// Return cached data
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(data); err != nil {
				logger.Errorf("Failed to encode cached data for %s: %v", key, err)
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
			return
		}
		
		// Fetch fresh data
		data, err := fetcher()
		if err != nil {
			logger.Errorf("Failed to fetch data for %s: %v", key, err)
			http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
			return
		}
		
		// Cache the data and get ETag
		etag := cache.Set(key, data, ttl)
		
		// Set response headers
		w.Header().Set("ETag", etag)
		w.Header().Set("Cache-Control", "private, max-age=300") // 5 minutes browser cache
		w.Header().Set("Content-Type", "application/json")
		
		// Return fresh data
		if err := json.NewEncoder(w).Encode(data); err != nil {
			logger.Errorf("Failed to encode fresh data for %s: %v", key, err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

// InvalidateCache helper to invalidate cache entries on updates
func InvalidateCache(cache interfaces.HTTPCacheService, pattern string) {
	cache.Invalidate(pattern)
}