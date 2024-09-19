package constants

type CacheType string // 缓存类型

var (
	MEMORY_CACHE CacheType = "MemoryCache" // 内存缓存
	REDIS_CACHE  CacheType = "RedisCache"  // Redis缓存
)
