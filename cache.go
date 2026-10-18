package godns

import "time"

type DnsCache struct{}

type DnsCacheKey struct {
	Qname  string
	Qtype  Type
	Qclass Class
}

type DnsCacheEntry struct {
	timestamp  time.Time
	ttl        int
	hits       int
	answers    []DnsRecord
	authority  []DnsRecord
	additional []DnsRecord
}

func NewDnsCacheEntry(answers []DnsRecord, authority []DnsRecord, additional []DnsRecord) {

}

func (c *DnsCache) Put(key DnsCacheKey, entry DnsCacheEntry) {}

// func (c *DnsCache) Get(key DnsCacheKey) (DnsCacheEntry, bool) {}
