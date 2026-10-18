package cache

type randomCache[K comparable, V any] struct {
	cacheMap map[K]*cacheEntry[V]
}

func (c *randomCache[K, V]) Put(key K, value V) {

}

func (c *randomCache[K, V]) Get(key K) (V, bool) {
	var value V
	return value, true
}

func (c *randomCache[K, V]) Evict(key K) {

}
