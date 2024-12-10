// persist/cache.go
package persist

import (
	"fmt"
	"github.com/guarzo/canifly/internal/services/interfaces"
	"os"
	"path/filepath"
	"time"

	"github.com/patrickmn/go-cache"
)

const DefaultExpiration = 30 * time.Minute
const cleanupInterval = 32 * time.Minute

type Cache struct {
	cache *cache.Cache
}

func NewCache() *Cache {
	return &Cache{
		cache: cache.New(DefaultExpiration, cleanupInterval),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	value, found := c.cache.Get(key)
	if !found {
		return nil, false
	}
	byteSlice, ok := value.([]byte)
	return byteSlice, ok
}

func (c *Cache) Set(key string, value []byte, expiration time.Duration) {
	c.cache.Set(key, value, expiration)
}

func (c *Cache) SaveToFile(filename string, logger interfaces.Logger) error {
	items := c.cache.Items()

	serializable := make(map[string]cacheItem, len(items))
	for k, v := range items {
		byteSlice, ok := v.Object.([]byte)
		if !ok {
			logger.Warnf("Skipping key %s as its value is not []byte", k)
			continue
		}
		serializable[k] = cacheItem{
			Value:      byteSlice,
			Expiration: time.Unix(0, v.Expiration),
		}
	}

	if err := writeJSONToFile(filename, serializable); err != nil {
		logger.WithError(err).Errorf("Failed to save cache to %s", filename)
		return err
	}

	logger.Debugf("Cache saved to %s", filename)
	return nil
}

func (c *Cache) LoadFromFile(filename string, logger interfaces.Logger) error {
	var serializable map[string]cacheItem
	if err := readJSONFromFile(filename, &serializable); err != nil {
		if os.IsNotExist(err) {
			logger.Infof("Cache file does not exist: %s", filename)
			return nil
		}
		logger.WithError(err).Errorf("Failed to load cache from %s", filename)
		return err
	}

	for k, item := range serializable {
		ttl := time.Until(item.Expiration)
		if ttl > 0 {
			c.cache.Set(k, item.Value, ttl)
		} else {
			logger.Infof("Skipping expired cache item: %s", k)
		}
	}
	logger.Debugf("Cache successfully loaded from file: %s", filename)
	return nil
}

type cacheItem struct {
	Value      []byte
	Expiration time.Time
}

// DataStore Cache Wrappers

func (ds *DataStore) GetFromCache(key string) ([]byte, bool) {
	return ds.apiCache.Get(key)
}

func (ds *DataStore) SetToCache(key string, value []byte, expiration time.Duration) {
	ds.apiCache.Set(key, value, expiration)
}

func (ds *DataStore) LoadApiCache() error {
	filename := ds.generateCacheDataFileName()
	return ds.apiCache.LoadFromFile(filename, ds.logger)
}

func (ds *DataStore) SaveApiCache() error {
	filename := ds.generateCacheDataFileName()
	return ds.apiCache.SaveToFile(filename, ds.logger)
}

func (ds *DataStore) generateCacheDataFileName() string {
	path, err := ds.getWritableCachePath()
	if err != nil {
		ds.logger.WithError(err).Error("Error retrieving writable data path for cache")
		return ""
	}
	return fmt.Sprintf("%s/cache.json", path)
}

func (ds *DataStore) getWritableCachePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve writeable directory: %w", err)
	}

	pathSuffix := os.Getenv("PATH_SUFFIX")
	cachePath := filepath.Join(configDir, "canifly", "apicache")
	if pathSuffix != "" {
		cachePath = filepath.Join(cachePath, pathSuffix)
	}

	if err := os.MkdirAll(cachePath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create writable apicache directory: %w", err)
	}

	return cachePath, nil
}
