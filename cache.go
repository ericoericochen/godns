package godns

type CacheKey interface {
	Hash() string
}

type Cache[K CacheKey, V any] interface {
	Put(K, V)
	Get(K)
	Evict(K)
}

type DnsCache struct {
}

// func (c *DnsCache)
