package cache

import (
	"sync"
	"time"
)

// Cache is a key-value cache that expires and evicts entries according to a TTL.
type Cache struct {
	entries map[string]entry
	now     func() time.Time
	mu      sync.RWMutex
}

type entry struct {
	value  interface{}
	expiry time.Time
}

func (e *entry) isExpired(now time.Time) bool { return now.After(e.expiry) }

// New creates a new cache which evicts expired entries every expiryInterval.
func New(expiryInterval time.Duration) *Cache {
	entries := make(map[string]entry)
	c := &Cache{entries: entries, now: time.Now}
	go func() {
		for {
			time.Sleep(expiryInterval)
			c.evictExpired()
		}
	}()
	return c
}

func (c *Cache) evictExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.now()
	for k, v := range c.entries {
		if v.isExpired(now) {
			delete(c.entries, k)
		}
	}
}

// Len returns the number of values in the cache. This includes entries that have expired, but are not yet evicted.
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}

// Get returns the cached value associated with key.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.entries[key]
	if !ok || v.isExpired(c.now()) {
		return nil, false
	}
	return v.value, true
}

// Set associates key with given value in the cache. The value is invalidated after ttl has passed.
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	expiry := c.now().Add(ttl)
	c.entries[key] = entry{value: value, expiry: expiry}
}
