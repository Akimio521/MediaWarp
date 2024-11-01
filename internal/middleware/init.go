package middleware

import (
	"MediaWarp/internal/cache"
)

var (
	cacheManager cache.CacheManager
)

func init() {
	cacheManager = cache.GetCacheManager()
}
