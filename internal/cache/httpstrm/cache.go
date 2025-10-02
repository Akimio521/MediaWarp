package httpstrm

import (
	"sync"
	"time"
)

type entry struct {
	url       string
	expiresAt time.Time
}

type Cache struct {
	mu    sync.RWMutex
	items map[string]entry
}

func New() *Cache {
	return &Cache{items: make(map[string]entry)}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	value, ok := c.items[key]
	c.mu.RUnlock()
	if !ok {
		return "", false
	}
	if time.Now().After(value.expiresAt) {
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		return "", false
	}
	return value.url, true
}

func (c *Cache) Set(key, url string, ttl time.Duration) {
	if ttl <= 0 {
		c.Delete(key)
		return
	}
	c.mu.Lock()
	c.items[key] = entry{url: url, expiresAt: time.Now().Add(ttl)}
	c.mu.Unlock()
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}
