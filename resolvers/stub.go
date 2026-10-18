package resolvers

import (
	"bytes"

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

func (r *StubResolver) Resolve(message *godns.DnsMessage) (*godns.DnsMessage, error) {
	// validate dns query message
	err := godns.ValidateQueryDnsMessage(message)
	if err != nil {
		return nil, err
	}

	// check if response is cached
	question := message.Questions[0]
	cacheKey := cache.DnsCacheKey{
		Qname:  question.Qname,
		Qtype:  question.Qtype,
		Qclass: question.Qclass,
	}
	cached, hit := r.cache.Get(cacheKey)
	if hit {
		res := cachedValueToResponseDnsMessage(message, cached)
		res.Header.RecursionAvailable = true
		return res, nil
	}

	// encode dns message to bytes
	var buffer bytes.Buffer
	encoder := godns.NewDnsMessageEncoder(message)
	encoder.Write(&buffer)

	// send the dns message
	res, err := godns.SendDnsMessage(r.server, godns.UDP_DNS_PORT, buffer.Bytes())
	if err != nil {
		return nil, err
	}

	// cache the response
	cacheValue := cache.DnsCacheValue{
		Answers:    res.Answers,
		Authority:  res.Authority,
		Additional: res.Additional,
	}

	r.cache.Put(cacheKey, cacheValue)
	return res, err
}
