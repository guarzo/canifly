// services/interfaces/cache_service.go
package interfaces

import "time"

type CacheService interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte, expiration time.Duration)
	LoadCache() error
	SaveCache() error
}
