// services/interfaces/cache_service.go
package interfaces

import "time"

type CacheService interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte, expiration time.Duration)
	LoadCache() error
	SaveCache() error
}

type CacheRepository interface {
	GetFromCache(key string) ([]byte, bool)
	SetToCache(key string, value []byte, expiration time.Duration)
	LoadApiCache() error
	SaveApiCache() error
}
