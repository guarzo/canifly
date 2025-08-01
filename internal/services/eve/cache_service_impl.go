package eve

import (
	"time"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

// cacheServiceImpl implements the CacheService interface using CacheRepository
type cacheServiceImpl struct {
	logger    interfaces.Logger
	cacheRepo interfaces.CacheRepository
}

// NewCacheService creates a new cache service
func NewCacheService(logger interfaces.Logger, cacheRepo interfaces.CacheRepository) interfaces.CacheService {
	return &cacheServiceImpl{
		logger:    logger,
		cacheRepo: cacheRepo,
	}
}

func (c *cacheServiceImpl) Get(key string) ([]byte, bool) {
	return c.cacheRepo.Get(key)
}

func (c *cacheServiceImpl) Set(key string, value []byte, expiration time.Duration) {
	c.cacheRepo.Set(key, value, expiration)
}

func (c *cacheServiceImpl) LoadCache() error {
	return c.cacheRepo.LoadApiCache()
}

func (c *cacheServiceImpl) SaveCache() error {
	return c.cacheRepo.SaveApiCache()
}