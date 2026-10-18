package cache

import (
	"fmt"

	"github.com/ericoericochen/godns"
)

type EvictionPolicy int

const (
	LRU EvictionPolicy = iota
	LFU
	Random
)

type CacheOptions struct {
	MaxSize int
	EvictionPolicy
}

type DnsCache interface {
	Put(DnsCacheKey, DnsCacheValue)
	Get(DnsCacheKey) (DnsCacheValue, bool)
}

// dns cache key
type DnsCacheKey struct {
	Qname  string
	Qtype  godns.Type
	Qclass godns.Class
}

// dns cache value
type DnsCacheValue struct {
	Answers    []*godns.DnsRecord
	Authority  []*godns.DnsRecord
	Additional []*godns.DnsRecord
}

func (cv *DnsCacheValue) TTL() int {
	records := append(cv.Answers, cv.Authority...)
	records = append(records, cv.Additional...)
	var ttl int
	for _, record := range records {
		ttl = int(record.TTL)
		break
	}
	return ttl
}

type dnsCache struct {
	cache cache[DnsCacheKey, DnsCacheValue]
}

func (c *dnsCache) Put(key DnsCacheKey, value DnsCacheValue) {
	ttl := value.TTL()
	fmt.Println("ttl: ", ttl)
	c.cache.Put(key, value, ttl)
}

func (c *dnsCache) Get(key DnsCacheKey) (DnsCacheValue, bool) {
	return c.cache.Get(key)
}

func NewDnsCache(opts *CacheOptions) DnsCache {
	var cache cache[DnsCacheKey, DnsCacheValue]
	if opts.EvictionPolicy == Random {
		cache = newRandomCache[DnsCacheKey, DnsCacheValue](opts.MaxSize)
	}

	return &dnsCache{cache: cache}
}
