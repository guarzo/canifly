// persist/cache.go
package cacheStore

import (
	"fmt"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services/interfaces"
	"os"
	"path/filepath"
	"time"

	"github.com/patrickmn/go-cache"
)

const DefaultExpiration = 30 * time.Minute
const cleanupInterval = 32 * time.Minute
const cacheFileName = "cache.json"

var _ interfaces.CacheRepository = (*CacheStore)(nil)

type CacheStore struct {
	cache  *cache.Cache
	logger interfaces.Logger
}

func NewCacheStore(l interfaces.Logger) *CacheStore {
	return &CacheStore{
		cache:  cache.New(DefaultExpiration, cleanupInterval),
		logger: l,
	}
}

func (c *CacheStore) Get(key string) ([]byte, bool) {
	value, found := c.cache.Get(key)
	if !found {
		return nil, false
	}
	byteSlice, ok := value.([]byte)
	return byteSlice, ok
}

func (c *CacheStore) Set(key string, value []byte, expiration time.Duration) {
	c.cache.Set(key, value, expiration)
}

func (c *CacheStore) saveToFile(filename string) error {
	items := c.cache.Items()

	serializable := make(map[string]cacheItem, len(items))
	for k, v := range items {
		byteSlice, ok := v.Object.([]byte)
		if !ok {
			c.logger.Warnf("Skipping key %s as its value is not []byte", k)
			continue
		}
		serializable[k] = cacheItem{
			Value:      byteSlice,
			Expiration: time.Unix(0, v.Expiration),
		}
	}

	if err := persist.OldSaveJson(filename, serializable); err != nil {
		c.logger.WithError(err).Errorf("Failed to save cache to %s", filename)
		return fmt.Errorf("failed to write JSON file: %v", err)
	}

	c.logger.Debugf("CacheStore saved to %s", filename)
	return nil
}

func (c *CacheStore) loadFromFile(filename string) error {
	var serializable map[string]cacheItem
	// Internally this calls OldReadJson (via OldReadJson)
	if err := persist.OldReadJson(filename, &serializable); err != nil {
		if os.IsNotExist(err) {
			c.logger.Infof("CacheStore file does not exist: %s", filename)
			return nil
		}
		c.logger.WithError(err).Errorf("Failed to load cache from %s", filename)
		return err
	}

	for k, item := range serializable {
		ttl := time.Until(item.Expiration)
		if ttl > 0 {
			c.cache.Set(k, item.Value, ttl)
		} else {
			c.logger.Infof("Skipping expired cache item: %s", k)
		}
	}
	c.logger.Debugf("CacheStore successfully loaded from file: %s", filename)
	return nil
}

type cacheItem struct {
	Value      []byte
	Expiration time.Time
}

func (c *CacheStore) LoadApiCache() error {
	filename, err := c.generateCacheDataFileName()
	if err != nil {
		return err
	}
	return c.loadFromFile(filename)
}

func (c *CacheStore) SaveApiCache() error {
	filename, err := c.generateCacheDataFileName()
	if err != nil {
		return err
	}
	return c.saveToFile(filename)
}

func (c *CacheStore) generateCacheDataFileName() (string, error) {
	configPath, err := persist.GetWriteableSubPath(persist.AccountDir)
	if err != nil {
		return "", err
	}

	return filepath.Join(configPath, cacheFileName), nil
}
