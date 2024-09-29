package constants

type CacheType string // 缓存类型

var (
	MEMORY_CACHE CacheType = "Memory" // 内存缓存
	REDIS_CACHE  CacheType = "Redis"  // Redis缓存
)
