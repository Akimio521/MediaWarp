package cache

import (
	"MediaWarp/constants"
	"MediaWarp/internal/cache/memory"
	"MediaWarp/internal/config"
	"time"
)

// 统一缓存接口
//
// 新增缓存需要满足设计逻辑
// 先通过spaceName获取缓存空间，再根据通过cacheKey缓存空间中的缓存
// 缓存每个缓存空间都是一个独立的缓存池，不存在锁
type Cache interface {
	UpdateCache(spaceName string, cacheKey string, cache any, cacheDuration time.Duration)
	GetCache(spaceName string, cacheKey string) (data any, ok bool)
}

var c Cache

func init() {
	if c == nil {
		switch config.Cache.Type {
		case constants.MEMORY_CACHE:
			c = &memory.MemoryCache{}
		default:
			panic("未知的缓存类型")
		}
	}
}

// 更新缓存
//
// spaceName: 缓存空间名称
// cacheKey: 缓存Key
// cache: 缓存数据
func Update(spaceName string, cacheKey string, cache any, cacheDuration time.Duration) {
	c.UpdateCache(spaceName, cacheKey, cache, cacheDuration)
}

// 获取缓存
//
// spaceName: 缓存空间名称
// cacheKey: 缓存Key
// data: 缓存数据
// ok: 是否存在
func Get(spaceName string, cacheKey string) (data any, ok bool) {
	return c.GetCache(spaceName, cacheKey)
}
