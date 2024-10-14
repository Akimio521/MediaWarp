package cache

import (
	"MediaWarp/constants"
	"MediaWarp/internal/cache/memory"
	"MediaWarp/internal/config"
)

var c CacheManager

// 获取缓存
func GetCacheManager() CacheManager {
	if c == nil {
		switch config.GetConfig().CacheType {
		case constants.MEMORY_CACHE:
			c = &memory.MemoryCache{}
		default:
			panic("未知的缓存类型")
		}
	}
	return c
}
