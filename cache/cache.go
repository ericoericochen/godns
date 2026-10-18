package cache

import "time"

type cache[K comparable, V any] interface {
	Put(K, V, int) // set k = v in cache with ttl
	Get(K) (V, bool)
	Evict(K)
}

type cacheEntry[V any] struct {
	timestamp time.Time
	ttl       int
	hits      int
	value     V
}
