package cacheutils

import (
	"MediaWarp/constants"
	"MediaWarp/internal/utils/cacheutils/memory"
)

var (
	availableCaches = map[constants.CacheType]Cache{
		constants.MEMORY_CACHE: &memory.MemoryCache{},
	}
)

// 获取缓存
func GetCache(cacheType constants.CacheType) Cache {
	return availableCaches[cacheType]
}
