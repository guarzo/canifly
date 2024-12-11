// persist/cache.go
package persist

import (
	"fmt"
	"github.com/guarzo/canifly/internal/services/interfaces"
	"os"
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

	if err := saveJSONToFile(filename, serializable); err != nil {
		logger.WithError(err).Errorf("Failed to save cache to %s", filename)
		return fmt.Errorf("failed to write JSON file: %v", err)
	}

	logger.Debugf("Cache saved to %s", filename)
	return nil
}

func (c *Cache) LoadFromFile(filename string, logger interfaces.Logger) error {
	var serializable map[string]cacheItem
	// Internally this calls readJSONFromFile (via readJSONFromFile)
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
	filename, err := ds.generateCacheDataFileName()
	if err != nil {
		return err
	}
	return ds.apiCache.LoadFromFile(filename, ds.logger)
}

func (ds *DataStore) SaveApiCache() error {
	filename, err := ds.generateCacheDataFileName()
	if err != nil {
		return err
	}
	return ds.apiCache.SaveToFile(filename, ds.logger)
}

func (ds *DataStore) generateCacheDataFileName() (string, error) {
	return getConfigFileName(cacheFileName)
}
