package resolvers

import (
	"github.com/ericoericochen/godns"
	"github.com/ericoericochen/godns/cache"
)

type StubResolver struct {
	server string
	cache  cache.DnsCache
}

func NewStubResolver(server string, cache cache.DnsCache) DnsResolver {
	resolver := &StubResolver{server: server, cache: cache}
	return resolver
}

func (r *StubResolver) Resolve(message godns.DnsMessage) godns.DnsMessage {
	// validate dns query message
	err := validateQueryDnsMessage(message)
	if err != nil {
		return formatErrorResponse(message)
	}

	// check if response is cached
	// question := message.Questions[0]
	// cacheKey := cache.DnsCacheKey{
	// 	Qname:  question.Qname,
	// 	Qtype:  question.Qtype,
	// 	Qclass: question.Qclass,
	// }
	// cached, hit := r.cache.Get(cacheKey)
	// if hit {
	// 	res := cachedValueToResponseDnsMessage(message, cached)
	// 	res.Header.RecursionAvailable = true
	// 	return res, nil
	// }

	// send the dns message
	res := SendDnsMessage(r.server, UDP_DNS_PORT, message)
	return res
}
