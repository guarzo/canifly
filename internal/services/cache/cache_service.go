package cache

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

type cacheEntry struct {
	data       interface{}
	etag       string
	expiration time.Time
}

type HTTPCacheService struct {
	entries map[string]*cacheEntry
	mu      sync.RWMutex
	logger  interfaces.Logger
}

// NewHTTPCacheService creates a new HTTP cache service
func NewHTTPCacheService(logger interfaces.Logger) *HTTPCacheService {
	service := &HTTPCacheService{
		entries: make(map[string]*cacheEntry),
		logger:  logger,
	}

	// Start cleanup goroutine
	go service.cleanupExpired()

	return service
}

// Get retrieves cached data and its ETag
func (s *HTTPCacheService) Get(key string) (interface{}, string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, exists := s.entries[key]
	if !exists {
		return nil, "", false
	}

	// Check if expired
	if time.Now().After(entry.expiration) {
		return nil, "", false
	}

	return entry.data, entry.etag, true
}

// Set stores data with a TTL and returns generated ETag
func (s *HTTPCacheService) Set(key string, data interface{}, ttl time.Duration) string {
	// Generate ETag from data
	etag := s.generateETag(data)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.entries[key] = &cacheEntry{
		data:       data,
		etag:       etag,
		expiration: time.Now().Add(ttl),
	}

	s.logger.Debugf("Cached %s with TTL %v", key, ttl)

	return etag
}

// Invalidate removes entries matching the pattern
func (s *HTTPCacheService) Invalidate(pattern string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	keysToDelete := []string{}
	for key := range s.entries {
		if strings.Contains(key, pattern) {
			keysToDelete = append(keysToDelete, key)
		}
	}

	for _, key := range keysToDelete {
		delete(s.entries, key)
		s.logger.Debugf("Invalidated cache key: %s", key)
	}
}

// Clear removes all cache entries
func (s *HTTPCacheService) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.entries = make(map[string]*cacheEntry)
	s.logger.Info("Cache cleared")
}

// generateETag creates an ETag from the data
func (s *HTTPCacheService) generateETag(data interface{}) string {
	// Convert data to JSON for consistent hashing
	jsonData, err := json.Marshal(data)
	if err != nil {
		// Fallback to timestamp-based ETag
		return hex.EncodeToString([]byte(time.Now().String()))
	}

	hash := md5.Sum(jsonData)
	return hex.EncodeToString(hash[:])
}

// cleanupExpired periodically removes expired entries
func (s *HTTPCacheService) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		keysToDelete := []string{}

		for key, entry := range s.entries {
			if now.After(entry.expiration) {
				keysToDelete = append(keysToDelete, key)
			}
		}

		for _, key := range keysToDelete {
			delete(s.entries, key)
		}

		if len(keysToDelete) > 0 {
			s.logger.Debugf("Cleaned up %d expired cache entries", len(keysToDelete))
		}
		s.mu.Unlock()
	}
}
