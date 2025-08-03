package interfaces

import "time"

// CacheableService handles cache operations
type CacheableService interface {
	SaveCache() error
	LoadCache() error
	SaveEsiCache() error
	Get(key string) ([]byte, bool)
	Set(key string, value []byte, expiration time.Duration)
}
