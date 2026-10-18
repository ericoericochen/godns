package cache

import "time"

type randomCache[K comparable, V any] struct {
	cacheMap map[K]*cacheEntry[V]
	maxSize  int
}

func newRandomCache[K comparable, V any](maxSize int) cache[K, V] {
	cache := &randomCache[K, V]{}
	cache.cacheMap = map[K]*cacheEntry[V]{}
	cache.maxSize = maxSize
	return cache
}

func (c *randomCache[K, V]) Put(key K, value V, ttl int) {
	entry := &cacheEntry[V]{
		timestamp: time.Now(),
		ttl:       ttl,
		hits:      0,
		value:     value,
	}
	c.cacheMap[key] = entry
}

func (c *randomCache[K, V]) Get(key K) (V, bool) {
	var value V
	entry, hit := c.cacheMap[key]
	if hit {
		value = entry.value
		entry.hits += 1
	}
	return value, hit
}

func (c *randomCache[K, V]) Evict(key K) {
	delete(c.cacheMap, key)
}
