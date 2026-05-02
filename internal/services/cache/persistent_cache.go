package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

// persistentCacheEntry represents a cached item with expiration
type persistentCacheEntry struct {
	data       []byte
	expiration time.Time
}

// PersistentCacheService is a disk-backed in-memory cache with TTL support.
type PersistentCacheService struct {
	logger   interfaces.Logger
	basePath string
	mu       sync.RWMutex
	entries  map[string]*persistentCacheEntry
	stopCh   chan struct{}
}

// NewPersistentCacheService creates a new persistent cache rooted at basePath
// and starts a background goroutine that periodically evicts expired entries.
func NewPersistentCacheService(basePath string, logger interfaces.Logger) *PersistentCacheService {
	c := &PersistentCacheService{
		logger:   logger,
		basePath: basePath,
		entries:  make(map[string]*persistentCacheEntry),
		stopCh:   make(chan struct{}),
	}
	go c.cleanupExpired()
	return c
}

// Get retrieves a cached value if present and not expired.
func (c *PersistentCacheService) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(entry.expiration) {
		c.mu.RUnlock()
		c.mu.Lock()
		delete(c.entries, key)
		c.mu.Unlock()
		c.mu.RLock()
		return nil, false
	}

	return entry.data, true
}

// Set stores a value with the given expiration TTL.
func (c *PersistentCacheService) Set(key string, value []byte, expiration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = &persistentCacheEntry{
		data:       value,
		expiration: time.Now().Add(expiration),
	}
}

// SaveCache flushes the cache contents to disk as api_cache.json.
func (c *PersistentCacheService) SaveCache() error {
	c.mu.RLock()
	cacheCopy := make(map[string][]byte, len(c.entries))
	for k, v := range c.entries {
		cacheCopy[k] = v.data
	}
	c.mu.RUnlock()

	if err := os.MkdirAll(c.basePath, 0o755); err != nil {
		return err
	}
	path := filepath.Join(c.basePath, "api_cache.json")
	data, err := json.Marshal(cacheCopy)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// LoadCache loads cache contents from disk if api_cache.json exists.
func (c *PersistentCacheService) LoadCache() error {
	path := filepath.Join(c.basePath, "api_cache.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	cache := make(map[string][]byte)
	if err := json.Unmarshal(data, &cache); err != nil {
		return err
	}
	c.mu.Lock()
	c.entries = make(map[string]*persistentCacheEntry, len(cache))
	for k, v := range cache {
		c.entries[k] = &persistentCacheEntry{
			data:       v,
			expiration: time.Now().Add(24 * time.Hour),
		}
	}
	c.mu.Unlock()
	return nil
}

// SaveEsiCache is an alias for SaveCache used by the ESI service interface.
func (c *PersistentCacheService) SaveEsiCache() error {
	return c.SaveCache()
}

// Clear removes all cache entries.
func (c *PersistentCacheService) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*persistentCacheEntry)
	c.logger.Info("EVE data cache cleared")
}

// Shutdown stops the background cleanup goroutine.
func (c *PersistentCacheService) Shutdown() {
	close(c.stopCh)
	c.logger.Info("Persistent cache service shutdown initiated")
}

// cleanupExpired periodically removes expired entries.
func (c *PersistentCacheService) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopCh:
			c.logger.Info("EVE data cache cleanup goroutine shutting down")
			return
		case <-ticker.C:
			c.mu.Lock()
			now := time.Now()
			keysToDelete := []string{}

			for key, entry := range c.entries {
				if now.After(entry.expiration) {
					keysToDelete = append(keysToDelete, key)
				}
			}

			for _, key := range keysToDelete {
				delete(c.entries, key)
			}

			if len(keysToDelete) > 0 {
				c.logger.Debugf("Cleaned up %d expired EVE data cache entries", len(keysToDelete))
			}
			c.mu.Unlock()
		}
	}
}

var _ interfaces.CacheableService = (*PersistentCacheService)(nil)
