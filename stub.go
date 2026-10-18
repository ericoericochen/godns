package godns

import (
	"fmt"
	"time"
)

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

func (r *StubResolver) Resolve(message DnsMessage) DnsMessage {
	// validate dns query message
	err := validateQueryDnsMessage(message)
	if err != nil {
		return toResponse(message, RcodeFormatError)
	}

	// build cache key
	question := message.Questions[0]
	cacheKey := DnsCacheKey{}
	cacheKey.Qname = question.Qname
	cacheKey.Qtype = question.Qtype
	cacheKey.Qclass = question.Qclass

	// check if response is cached
	entry, hit := r.cache.Get(cacheKey)
	if hit {
		fmt.Println("cache hit")
		res := NewDnsMessageBuilder().
			From(message).
			Response().
			RecursionAvailable(true).
			Answers(entry.Answers).
			Authority(entry.Authority).
			Additional(entry.Additional).
			Build()

		return res
	}

	// send the dns message to recursive resolver
	res := SendDnsMessage(r.Server, UDP_DNS_PORT, message)

	// cache the response if successful response
	if r.useCache() && res.Header.Rcode == RcodeNoError {
		fmt.Println("caching response")
		entry := DnsCacheEntry{
			Timestamp:  time.Now(),
			TTL:        300,
			Answers:    res.Answers,
			Authority:  res.Authority,
			Additional: res.Additional,
		}
		r.cache.Put(cacheKey, entry)
	}

	return res
}

func (r *StubResolver) useCache() bool {
	return r.cache != nil
}
