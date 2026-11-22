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

func NewStubResolver(server netip.Addr, cache *DnsCache) *StubResolver {
	if !server.Is4() {
		panic(fmt.Sprintf("%v is not an ipv4 address", server))
	}
	if cache == nil {
		panic("Cache must not be nil")
	}
	resolver := &StubResolver{server: server, cache: cache}
	return resolver
}

func (r *StubResolver) Resolve(query DnsMessage) DnsMessage {
	// validate dns query message
	err := validateQueryDnsMessage(query)
	if err != nil {
		return query.SetQR(true).SetRcode(RcodeFormatError)
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
	r.cacheResponse(query, res)

	return res
}

func (r *StubResolver) getCached(query DnsMessage) *DnsMessage {
	return getCachedResponseFromDnsCache(r.cache, query)
}

func (r *StubResolver) cacheResponse(query DnsMessage, response DnsMessage) {
	if response.Header.Rcode == RcodeNoError {
		saveResponseToDnsCache(r.cache, query, response)
	}
}
