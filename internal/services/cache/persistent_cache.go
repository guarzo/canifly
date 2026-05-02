package cache

import (
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
	logger     interfaces.Logger
	storage    interfaces.StorageService
	mu         sync.RWMutex
	entries    map[string]*persistentCacheEntry
	stopCh     chan struct{}
	stopOnce   sync.Once
	wg         sync.WaitGroup
}

// NewPersistentCacheService creates a new persistent cache that delegates disk
// IO to the supplied storage service and starts a background goroutine that
// periodically evicts expired entries.
func NewPersistentCacheService(storage interfaces.StorageService, logger interfaces.Logger) *PersistentCacheService {
	c := &PersistentCacheService{
		logger:  logger,
		storage: storage,
		entries: make(map[string]*persistentCacheEntry),
		stopCh:  make(chan struct{}),
	}
	c.wg.Add(1)
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

// SaveCache flushes the cache contents to disk via the storage service.
func (c *PersistentCacheService) SaveCache() error {
	c.mu.RLock()
	cacheCopy := make(map[string][]byte, len(c.entries))
	for k, v := range c.entries {
		cacheCopy[k] = v.data
	}
	c.mu.RUnlock()

	return c.storage.SaveAPICache(cacheCopy)
}

// LoadCache loads cache contents from disk via the storage service.
func (c *PersistentCacheService) LoadCache() error {
	cache, err := c.storage.LoadAPICache()
	if err != nil {
		return err
	}
	c.mu.Lock()
	c.entries = make(map[string]*persistentCacheEntry, len(cache))
	for k, v := range cache {
		// On-disk format stores only []byte values, not original expirations.
		// Reloaded entries are given a fresh TTL; this matches prior behavior.
		c.entries[k] = &persistentCacheEntry{
			data:       v,
			expiration: time.Now().Add(24 * time.Hour),
		}
	}
	c.mu.Unlock()
	return nil
}

// SaveEsiCache exists to satisfy the legacy ESIService interface contract.
// TODO: drop once that interface no longer requires it.
func (c *PersistentCacheService) SaveEsiCache() error {
	return c.SaveCache()
}

// Clear removes all cache entries.
func (c *PersistentCacheService) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*persistentCacheEntry)
	c.logger.Info("persistent cache cleared")
}

// Shutdown stops the background cleanup goroutine and blocks until it exits.
// Safe to call multiple times.
func (c *PersistentCacheService) Shutdown() {
	c.stopOnce.Do(func() {
		close(c.stopCh)
		c.logger.Info("persistent cache shutdown initiated")
	})
	c.wg.Wait()
}

// cleanupExpired periodically removes expired entries.
func (c *PersistentCacheService) cleanupExpired() {
	defer c.wg.Done()
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopCh:
			c.logger.Info("persistent cache cleanup goroutine shutting down")
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
				c.logger.Debugf("persistent cache cleaned up %d expired entries", len(keysToDelete))
			}
			c.mu.Unlock()
		}
	}
}

var _ interfaces.CacheableService = (*PersistentCacheService)(nil)
