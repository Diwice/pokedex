package cache

import (
	"sync"
	"time"
)

type Cache struct {
	mu     sync.Mutex
	caches map[string]cache_entry
}

type cache_entry struct {
	val        []byte
	created_at time.Time
}

func New_Cache(interval time.Duration) Cache {
	new_cache := Cache{
		caches: map[string]cache_entry{},
	}
	go new_cache.reap_loop(interval)

	return new_cache
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()

	c.caches[key] = cache_entry{
		val: value,
		created_at: time.Now(),
	}

	c.mu.Unlock()
} 

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.caches[key]
	if !ok {
		return []byte{}, false
	}
	
	return entry.val, true
}

func (c *Cache) reap_loop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	
	for range ticker.C{
		c.mu.Lock()

		for i, v := range c.caches {
			if time.Since(v.created_at) >= interval {
				delete(c.caches, i)
			}
		}

		c.mu.Unlock()
	}
} 
