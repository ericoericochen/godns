package godns

// Create dns cache key from dns query message
func queryToCacheKey(message DnsMessage) DnsCacheKey {
	qname := message.Qname()
	qtype := message.Qtype()
	qclass := message.Qclass()
	return DnsCacheKey{
		Qname:  qname,
		Qtype:  qtype,
		Qclass: qclass,
	}
}

// Get cached response from dns cache
func getCachedResponseFromDnsCache(cache *DnsCache, query DnsMessage) *DnsMessage {
	cacheKey := queryToCacheKey(query)
	entry := cache.Get(cacheKey)
	if entry == nil {
		return nil
	}

	answers := entry.Value.Answers
	authority := entry.Value.Authority
	additional := entry.Value.Additional

	res := query.
		SetQR(true).
		SetRecursionAvailable(true).
		SetAnswers(answers).
		SetAuthority(authority).
		SetAdditional(additional)

	return &res
}

// Save response in cache
func saveResponseToDnsCache(cache *DnsCache, query DnsMessage, response DnsMessage) {
	cacheKey := queryToCacheKey(query)
	cacheValue := DnsCacheValue{
		Answers:    response.Answers,
		Authority:  response.Authority,
		Additional: response.Additional,
	}
	cache.Put(cacheKey, cacheValue)
}
