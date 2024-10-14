package middleware

import (
	"MediaWarp/internal/cache"
	"MediaWarp/internal/config"
	"MediaWarp/internal/log"
)

var (
	cfg          *config.ConfigManager
	logger       *log.LoggerManager
	cacheManager cache.CacheManager
)

func init() {
	cfg = config.GetConfig()
	logger = log.GetLogger()
	cacheManager = cache.GetCacheManager()
}
