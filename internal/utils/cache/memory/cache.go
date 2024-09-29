package memory

import (
	"sync"
	"time"
)

// 内存缓存
type MemoryCache struct {
	cacheSpaceMap sync.Map
}

// 获取子缓存空间
//
// 通过spaceName获取子缓存空间，若cacheSpaceMap中无该子空间，则实例化一个新的子空间后返回
func (memoryCache *MemoryCache) getCacheSpace(spaceName string) *MemoryCacheSpace {
	if cacheSpace, ok := memoryCache.cacheSpaceMap.Load(spaceName); ok { // 需要获取的子缓存空间已被创建
		return cacheSpace.(*MemoryCacheSpace)
	}

	// 生成一个新的子缓存空间
	newCacheSpace := &MemoryCacheSpace{}

	// fmt.Println("创建子缓存空间", spaceName)
	memoryCache.cacheSpaceMap.Store(spaceName, newCacheSpace)
	return newCacheSpace
}

// 更新对应子缓存空间中的数据
func (memoryCache *MemoryCache) UpdateCache(spaceName string, cacheKey string, cacheData any, cacheDuration time.Duration) {
	memoryCache.getCacheSpace(spaceName).UpdateCache(cacheKey, cacheData, cacheDuration)
}

// 从对应子缓存空间中获取数据
func (memoryCache *MemoryCache) GetCache(spaceName string, cacheKey string) (any, bool) {
	return memoryCache.getCacheSpace(spaceName).GetCache(cacheKey)
}
