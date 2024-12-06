// persist/cache.go
package persist

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/guarzo/canifly/internal/utils/xlog"
)

const DefaultExpiration = 30 * time.Minute
const cleanupInterval = 32 * time.Minute

var ApiCache *Cache

func GetWritableCachePath() (string, error) {
	configDir, err := os.UserConfigDir()

	if err != nil {
		return "", fmt.Errorf("failed to retrieve writeable directory: %w", err)
	}

	pathSuffix := os.Getenv("PATH_SUFFIX")
	cachePath := filepath.Join(configDir, "canifly", "apicache")
	if pathSuffix != "" {
		cachePath = filepath.Join(configDir, "canifly", "apicache", pathSuffix)
	}

	if err = os.MkdirAll(cachePath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create writable apicache directory: %w", err)
	}

	return cachePath, nil
}

// Cache struct to manage in-memory caching with optional persistence.
type Cache struct {
	cache *cache.Cache
}

func NewCache() *Cache {
	return &Cache{
		cache: cache.New(DefaultExpiration, cleanupInterval),
	}
}

// Get retrieves a value from the cache by key.
func (c *Cache) Get(key string) ([]byte, bool) {
	value, found := c.cache.Get(key)
	if !found {
		return nil, false
	}
	byteSlice, ok := value.([]byte)
	return byteSlice, ok
}

// Set stores data in the cache with the associated key and expiration duration.
func (c *Cache) Set(key string, value []byte, expiration time.Duration) {
	c.cache.Set(key, value, expiration)
}

// SaveToFile saves the entire cache to a file in JSON format.
func (c *Cache) SaveToFile() error {
	items := c.cache.Items()

	// Convert to map[string]cacheItem for serialization
	serializable := make(map[string]cacheItem, len(items))
	for k, v := range items {
		byteSlice, ok := v.Object.([]byte)
		if !ok {
			xlog.Logf("Skipping key %s as its value is not []byte", k)
			continue
		}
		serializable[k] = cacheItem{
			Value:      byteSlice,
			Expiration: time.Unix(0, v.Expiration),
		}
	}
	return WriteJSONToFile(GenerateCacheDataFileName(), serializable)
}

// LoadFromFile loads the cache from a JSON file.
func (c *Cache) LoadFromFile() error {
	var serializable map[string]cacheItem
	if err := ReadJSONFromFile(GenerateCacheDataFileName(), &serializable); err != nil {
		if os.IsNotExist(err) {
			xlog.Logf("Cache file does not exist: %s", GenerateCacheDataFileName())
			return nil
		}
		return err
	}

	// Set each item in the cache with its expiration check
	for k, item := range serializable {
		ttl := time.Until(item.Expiration)
		if ttl > 0 {
			c.cache.Set(k, item.Value, ttl)
		} else {
			xlog.Logf("Skipping expired cache item: %s", k)
		}
	}
	xlog.Logf("Cache successfully loaded from file: %s", GenerateCacheDataFileName())
	return nil
}

// GenerateCacheDataFileName generates the filename for storing cache data.
func GenerateCacheDataFileName() string {
	path, err := GetWritableCachePath()
	if err != nil {
		xlog.Logf("Error retrieving writable data path: %v\n", err)
		return ""
	}

	return fmt.Sprintf("%s/cache.json", path)
}

// cacheItem represents an item stored in the cache with its expiration.
type cacheItem struct {
	Value      []byte
	Expiration time.Time
}

func init() {
	ApiCache = NewCache()
	if ApiCache == nil {
		log.Fatal("Failed to initialize cache")
	}

	if err := ApiCache.LoadFromFile(); err != nil {
		xlog.Logf("Failed to load cache from file")
	}
}
