package cache

import "time"

// 统一缓存接口
//
// 新增缓存需要满足设计逻辑
// 先通过spaceName获取缓存空间，再根据通过cacheKey缓存空间中的缓存
// 缓存每个缓存空间都是一个独立的缓存池，不存在锁
type CacheManager interface {
	UpdateCache(spaceName string, cacheKey string, cache any, cacheDuration time.Duration)
	GetCache(spaceName string, cacheKey string) (data any, ok bool)
}
