package godns

import (
	"time"
)

type DnsCache struct {
	keyToEntries map[DnsCacheKey]*DnsCacheEntry
}

type DnsCacheKey struct {
	Qname  string
	Qtype  Type
	Qclass Class
}

type DnsCacheEntry struct {
	Timestamp  time.Time
	TTL        int
	Answers    []DnsRecord
	Authority  []DnsRecord
	Additional []DnsRecord
	hits       int
}

func NewDnsCache(capacity int) *DnsCache {
	cache := &DnsCache{}
	cache.keyToEntries = map[DnsCacheKey]*DnsCacheEntry{}
	return cache
}

func (c *DnsCache) Put(key DnsCacheKey, entry DnsCacheEntry) {
	c.keyToEntries[key] = &entry
}

func (c *DnsCache) Get(key DnsCacheKey) (DnsCacheEntry, bool) {
	var entry DnsCacheEntry
	ptr, hit := c.keyToEntries[key]
	if !hit {
		return entry, hit
	}
	return *ptr, hit
}
