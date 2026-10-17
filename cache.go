package godns

import (
	"fmt"
	"time"
)

// interfaces
type ICacheEntry[T any] interface {
	Timestamp() time.Time // timestamp of insertion
	TTL() int32           // time to live
	Value() T             // value of the cached data
	Hits() uint32         // number of puts / gets for this cache entry
	hit()                 // increment number of hits
}

type Cache[K, V any] interface {
	Put(K, V, *int32)
	Get(K) (ICacheEntry[V], bool)
	Evict(K)
}

// cache implementation

// dns cache key = (qname, qtype, qclass)
type DnsCacheKey struct {
	Qname  string
	Qtype  Type
	Qclass Class
}

func (k *DnsCacheKey) Hash() string {
	return fmt.Sprintf("%v:%v:%v", k.Qname, typeNameMap[k.Qtype], classNameMap[k.Qclass])
}

// dns cache value and entry

type DnsCacheValue struct {
	Answers    []*DNSRecord
	Authority  []*DNSRecord
	Additional []*DNSRecord
}

type CacheEntry[V any] struct {
	timestamp time.Time // when the cache entry was added to the cache
	ttl       int32
	hits      int32
	value     V
}

func (e *CacheEntry[V]) Value() V {
	return e.value
}

func (e *CacheEntry[V]) TTL() int32 {
	return e.ttl
}

func (e *CacheEntry[V]) Timestamp() time.Time {
	return e.timestamp
}

func (e *CacheEntry[V]) Hits() int32 {
	return e.hits
}

func (e *CacheEntry[V]) hit() {
	e.hits += 1
}

// dns cache - ttl cache with some eviction policy
type DnsCache struct {
}

type randomCache[K comparable, V any] struct {
	keyToEntry map[K]*CacheEntry[V]
}

func (c *randomCache[K, V]) Put(key K, value V, ttl *int32) {
	var ttlValue int32
	if ttl == nil {
		ttlValue = -1
	} else {
		ttlValue = *ttl
	}
	entry := &CacheEntry[V]{
		timestamp: time.Now(),
		ttl:       ttlValue,
		hits:      1,
		value:     value,
	}
	c.keyToEntry[key] = entry
}

func (c *randomCache[K, V]) Get(key K) (*CacheEntry[V], bool) {
	entry, exists := c.keyToEntry[key]
	return entry, exists
}

func (c *randomCache[K, V]) Evict(key K) {
	delete(c.keyToEntry, key)
}
