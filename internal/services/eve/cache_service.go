// services/cache/cache_service.go
package eve

import (
	"time"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

var _ interfaces.CacheService = (*cacheService)(nil)

type cacheService struct {
	logger       interfaces.Logger
	persistCache interfaces.CacheRepository // We'll define CacheRepository below
}

func NewCacheService(logger interfaces.Logger, persistCache interfaces.CacheRepository) interfaces.CacheService {
	return &cacheService{
		logger:       logger,
		persistCache: persistCache,
	}
}

func (c *cacheService) Get(key string) ([]byte, bool) {
	return c.persistCache.Get(key)
}

func (c *cacheService) Set(key string, value []byte, expiration time.Duration) {
	c.persistCache.Set(key, value, expiration)
}

func (c *cacheService) LoadCache() error {
	return c.persistCache.LoadApiCache()
}

func (c *cacheService) SaveCache() error {
	return c.persistCache.SaveApiCache()
}
