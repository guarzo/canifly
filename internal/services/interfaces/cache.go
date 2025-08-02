package interfaces

import "time"

// HTTPCacheService provides in-memory caching with ETag support for HTTP responses
type HTTPCacheService interface {
	// Get retrieves cached data and its ETag
	Get(key string) (data interface{}, etag string, found bool)
	
	// Set stores data with a TTL and returns generated ETag
	Set(key string, data interface{}, ttl time.Duration) string
	
	// Invalidate removes entries matching the pattern
	Invalidate(pattern string)
	
	// Clear removes all cache entries
	Clear()
}