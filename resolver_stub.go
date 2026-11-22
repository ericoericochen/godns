package godns

import (
	"fmt"
	"log"
	"net/netip"
)

type StubResolver struct {
	server netip.Addr
	cache  *DnsCache
}

const CLOUDFLARE_DNS_SERVER = "1.1.1.1"

func NewStubResolver() *StubResolver {
	ip, _ := netip.ParseAddr(CLOUDFLARE_DNS_SERVER)
	resolver := &StubResolver{server: ip}
	return resolver
}

func (r *StubResolver) WithServer(server netip.Addr) *StubResolver {
	r.server = server
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
		log.Println(fmt.Sprintf("cache hit for qname=%v, qtype=%v", query.Qname(), typeNameMap[query.Qtype()]))
		return *cachedMessage
	}

	// send the dns message to recursive resolver
	res := SendDnsMessage(r.server, UDP_DNS_PORT, query)

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

func (r *StubResolver) maybeCache(query DnsMessage, response DnsMessage) {
	if r.useCache() && response.Header.Rcode == RcodeNoError {
		cacheKey := queryToCacheKey(query)
		value := DnsCacheValue{
			Answers:    response.Answers,
			Authority:  response.Authority,
			Additional: response.Additional,
		}
		r.cache.Put(cacheKey, value)
	}
}

func (r *StubResolver) useCache() bool {
	return r.cache != nil
}
