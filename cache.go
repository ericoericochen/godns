package godns

import (
	"fmt"
	"time"
)

// interfaces
type CacheKey interface {
	Hash() string
}

type CacheEntry[T any] interface {
	Value() T     // value of the cached data
	Hits() uint32 // number of puts / gets for this cache entry
	hit()         // increment number of hits
}

type Cache[K CacheKey, V any] interface {
	Put(K, V)
	Get(K) (V, bool)
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

type DnsCacheEntry struct {
	timestamp time.Time // when the cache entry was added to the cache
	ttl       int32
	hits      int32
	data      *DnsCacheValue
}

func (e *DnsCacheEntry) TTL() int32 {
	return e.ttl
}

func (e *DnsCacheEntry) Timestamp() time.Time {
	return e.timestamp
}

func (e *DnsCacheEntry) Hits() int32 {
	return e.hits
}

func (e *DnsCacheEntry) hit() {
	e.hits += 1
}

// dns cache - ttl cache with some eviction policy
type DnsCache struct {
}

// TTL cache
type ttlCacheEntry[T any] interface {
	CacheEntry[T]
	TTL() int32
	Timestamp() time.Time
}

type ttlCache struct{}

// cache implementations with different eviction policies
type lruCache struct{}

type lfuCache struct{}

type randomCache[K CacheKey, V any] struct {
	// m map[]
}
