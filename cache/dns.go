package cache

import "github.com/ericoericochen/godns"

type EvictionPolicy int

const (
	LRU EvictionPolicy = iota
	LFU
	Random
)

type DnsCacheOptions struct {
	MaxSize int
	EvictionPolicy
}

type DnsCacheKey struct {
	Qname  string
	Qtype  godns.Type
	Qclass godns.Class
}

type DnsCacheValue struct {
	Answers    []*godns.DnsRecord
	Authority  []*godns.DnsRecord
	Additional []*godns.DnsRecord
}

func New(opts *DnsCacheOptions) Cache[DnsCacheKey, DnsCacheValue] {
	return nil
}
