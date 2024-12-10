// services/cache/cache_service.go
package cache

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

type cacheService struct {
	logger       *logrus.Logger
	persistCache interfaces.CacheRepository // We'll define CacheRepository below
}

func NewCacheService(logger *logrus.Logger, persistCache interfaces.CacheRepository) interfaces.CacheService {
	return &cacheService{
		logger:       logger,
		persistCache: persistCache,
	}
}

func (c *cacheService) Get(key string) ([]byte, bool) {
	return c.persistCache.GetFromCache(key)
}

func (c *cacheService) Set(key string, value []byte, expiration time.Duration) {
	c.persistCache.SetToCache(key, value, expiration)
}

func (c *cacheService) LoadCache() error {
	return c.persistCache.LoadApiCache()
}

func (c *cacheService) SaveCache() error {
	return c.persistCache.SaveApiCache()
}
