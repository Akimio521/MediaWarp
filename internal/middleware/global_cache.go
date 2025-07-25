package middleware

import (
	"sync"
	"time"
)

var (
	instance *TTLCache
	once     sync.Once
)

type TTLCache struct {
	items map[string]*cacheItem
	mu    sync.RWMutex
}

type cacheItem struct {
	value      interface{}
	expiration int64 // UnixNano 时间戳
}

// 获取全局缓存实例（懒加载）
func GetCache() *TTLCache {
	once.Do(func() {
		instance = &TTLCache{
			items: make(map[string]*cacheItem),
		}
		go instance.cleanupExpired(1 * time.Minute) // 后台清理协程
	})
	return instance
}

func (c *TTLCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = &cacheItem{
		value:      value,
		expiration: time.Now().Add(ttl).UnixNano(),
	}
}

func (c *TTLCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()

	if !ok || time.Now().UnixNano() > item.expiration {
		return nil, false
	}
	return item.value, true
}

// 后台清理协程（无需 stop，程序退出时自动终止）
func (c *TTLCache) cleanupExpired(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for key, item := range c.items {
			if time.Now().UnixNano() > item.expiration {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}
