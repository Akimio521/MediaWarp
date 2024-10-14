package memory

import (
	"sync"
	"time"
)

// 内存缓存子缓存空间
type MemoryCacheSpace struct {
	cachePool sync.Map
}

// 更新缓存中的数据
func (memoryCacheSpace *MemoryCacheSpace) UpdateCache(cacheKey string, cacheData any, cacheDuration time.Duration) {
	cacheItem := &CacheItem{
		Data:   cacheData,
		Expiry: time.Now().Add(cacheDuration),
	}

	// fmt.Println("更新缓存")
	memoryCacheSpace.cachePool.Store(cacheKey, cacheItem)
}

// 从缓存中获取数据
func (memoryCacheSpace *MemoryCacheSpace) GetCache(cacheKey string) (any, bool) {
	cacheIF, ok := memoryCacheSpace.cachePool.Load(cacheKey)

	if ok { // 找到缓存
		cacheItem := cacheIF.(*CacheItem)
		if time.Now().After(cacheItem.Expiry) { //缓存已过期
			memoryCacheSpace.cachePool.Delete(cacheKey) // 删除缓存
			return nil, false                           // 修改为未找到缓存
		}

		// fmt.Println("命中缓存")
		return cacheItem.Data, true
	}

	// 未找到缓存
	return nil, false
}
