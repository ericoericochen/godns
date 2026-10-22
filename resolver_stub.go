package godns

type StubResolver struct {
	Server string
	cache  *DnsCache
}

const CLOUDFLARE_DNS_SERVER = "1.1.1.1"

func NewStubResolver() *StubResolver {
	resolver := &StubResolver{Server: CLOUDFLARE_DNS_SERVER}
	return resolver
}

func (r *StubResolver) WithServer(server string) *StubResolver {
	r.Server = server
	return r
}

func (r *StubResolver) WithCache(cache *DnsCache) *StubResolver {
	r.cache = cache
	return r
}

func (r *StubResolver) Resolve(query DnsMessage) DnsMessage {
	// validate dns query message
	err := validateQueryDnsMessage(query)
	if err != nil {
		return toResponse(query, RcodeFormatError)
	}

	// try to get cached response
	cachedMessage := r.getCached(query)
	if cachedMessage != nil {
		return *cachedMessage
	}

	// send the dns message to recursive resolver
	res := SendDnsMessage(r.Server, UDP_DNS_PORT, query)

	// cache the response if successful response
	r.maybeCache(query, res)

	return res
}

func (r *StubResolver) getCached(query DnsMessage) *DnsMessage {
	if !r.useCache() {
		return nil
	}
	cacheKey := queryToCacheKey(query)
	entry := r.cache.Get(cacheKey)
	if entry == nil {
		return nil
	}

	res := NewDnsMessageBuilder().
		From(query).
		Response().
		RecursionAvailable(true).
		Answers(entry.Value.Answers).
		Authority(entry.Value.Authority).
		Additional(entry.Value.Additional).
		Build()

	return &res
}

func (r *StubResolver) maybeCache(query DnsMessage, response DnsMessage) {
	if r.useCache() && response.Header.Rcode == RcodeNoError {
		cacheKey := queryToCacheKey(query)
		value := DnsCacheValue{
			Answers:    response.Answers,
			Authority:  response.Authority,
			Additional: response.Additional,
			TTL:        300,
		}
		r.cache.Put(cacheKey, value)
	}
}

func (r *StubResolver) useCache() bool {
	return r.cache != nil
}

// utilties

// create dns cache key from dns query message
func queryToCacheKey(message DnsMessage) DnsCacheKey {
	question := message.Questions[0]
	cacheKey := DnsCacheKey{}
	cacheKey.Qname = question.Qname
	cacheKey.Qtype = question.Qtype
	cacheKey.Qclass = question.Qclass
	return cacheKey
}
