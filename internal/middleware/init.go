package middleware

import (
	"MediaWarp/internal/cache"
	"MediaWarp/internal/log"
)

var (
	logger       *log.LoggerManager
	cacheManager cache.CacheManager
)

func init() {
	logger = log.GetLogger()
	cacheManager = cache.GetCacheManager()
}
