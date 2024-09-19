package cachepkg

import (
	"MediaWarp/constants"
	"MediaWarp/pkg/cachepkg/memory"
)

var (
	AvailableCaches map[constants.CacheType]Cache
)

func init() {
	// 初始化 AvailableCaches
	AvailableCaches = make(map[constants.CacheType]Cache)
	AvailableCaches[constants.MEMORY_CACHE] = &memory.MemoryCache{}
}
