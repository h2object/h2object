package object

import (
	"time"	
	"sync"
	"github.com/pmylund/go-cache"
)

type MemoryCache struct{
	sync.RWMutex
	base *cache.Cache
}

func NewMutexCache(defaultExpiration time.Duration, cleanupInterval time.Duration) *MemoryCache {
	return &MemoryCache{
		base: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (c *MemoryCache) Set(k string, v interface{}, exp time.Duration) {
	c.Lock()
	defer c.Unlock()
	c.base.Set(k, v, exp)
}

func (c *MemoryCache) Add(k string, v interface{}, exp time.Duration) error {
	c.Lock()
	defer c.Unlock()
	return c.base.Add(k, v, exp)	
}

func (c *MemoryCache) Get(k string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()
	return c.base.Get(k)	
}

func (c *MemoryCache) Delete(k string) {
	c.Lock()
	defer c.Unlock()
	c.base.Delete(k)
}

func (c *MemoryCache) DeleteExpired() {
	c.Lock()
	defer c.Unlock()
	c.base.DeleteExpired()
}