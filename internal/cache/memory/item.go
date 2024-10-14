package memory

import "time"

// 存入缓存池的数据结构
type CacheItem struct {
	Data   any       // 具体数据
	Expiry time.Time // 数据过期时间
}
