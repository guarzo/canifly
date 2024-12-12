package eve

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services/interfaces"
	"github.com/patrickmn/go-cache"
)

const (
	DefaultExpiration = 30 * time.Minute
	cleanupInterval   = 32 * time.Minute
	cacheFileName     = "cache.json"
)

var _ interfaces.CacheRepository = (*CacheStore)(nil)

type CacheStore struct {
	cache    *cache.Cache
	logger   interfaces.Logger
	fs       persist.FileSystem
	basePath string
}

func NewCacheStore(logger interfaces.Logger, fs persist.FileSystem, basePath string) *CacheStore {
	return &CacheStore{
		cache:    cache.New(DefaultExpiration, cleanupInterval),
		logger:   logger,
		fs:       fs,
		basePath: basePath,
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

func (c *CacheStore) LoadApiCache() error {
	filename := filepath.Join(c.basePath, cacheFileName)
	var serializable map[string]cacheItem

	if _, err := c.fs.Stat(filename); os.IsNotExist(err) {
		c.logger.Infof("CacheStore file does not exist: %s", filename)
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to stat cache file: %w", err)
	}

	if err := persist.ReadJsonFromFile(c.fs, filename, &serializable); err != nil {
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

func (c *CacheStore) SaveApiCache() error {
	filename := filepath.Join(c.basePath, cacheFileName)
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

	if err := persist.SaveJsonToFile(c.fs, filename, serializable); err != nil {
		c.logger.WithError(err).Errorf("Failed to save cache to %s", filename)
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	c.logger.Debugf("CacheStore saved to %s", filename)
	return nil
}

type cacheItem struct {
	Value      []byte
	Expiration time.Time
}

func (c *CacheStore) CacheItemsForTest() map[string]cache.Item {
	return c.cache.Items()
}

func (c *CacheStore) CacheSetForTest(key string, value interface{}) {
	c.cache.Set(key, value, DefaultExpiration)
}
