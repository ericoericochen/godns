package godns

import (
	"fmt"
	"log"
	"strings"
)

type RecursiveResolver struct {
	cache *DnsCache
}

func NewRecursiveResolver() *RecursiveResolver {
	return &RecursiveResolver{}
}

func (r *RecursiveResolver) WithCache(cache *DnsCache) *RecursiveResolver {
	r.cache = cache
	return r
}

func (r *RecursiveResolver) Resolve(query DnsMessage) DnsMessage {
	// validate dns message
	return r.resolve(query)
}

func (r *RecursiveResolver) resolve(query DnsMessage) DnsMessage {
	// check if query response is cached
	res, cached := r.getCachedResponse(query)
	if cached {
		return *res
	}

	// find the best starting nameserver to resolve against
	nameservers := r.getBestStartingNameserverIpv4s(query.Qname())

	for {
		nameserver := nameservers[0]
		log.Println(fmt.Sprintf("resolving qname=%v qtype=%v qclass=%v with ns %v",
			query.Qname(),
			typeNameMap[query.Qtype()],
			classNameMap[query.Qclass()],
			nameserver))
		res := SendDnsMessage(nameserver, UDP_DNS_PORT, query)
		if res.Header.Rcode != RcodeNoError {
			return res
		}

		// did we receive answer to query vs. referral
		if len(res.Answers) > 0 {
			r.cacheResponse(query, res)
			return res
		}

		// cache the referral nameservers and their ip addresses
		nextNameservers := r.getNextNameservers(res)
		nameservers = nextNameservers
	}
}

// response doesn't have answer and refers to next nameservers to query
func (r *RecursiveResolver) getNextNameservers(response DnsMessage) []string {
	// get nameserver records and their A records
	nameserverRecords := filterRecordsByType(response.Authority, TypeNS)
	nameserverIpv4Records := filterRecordsByType(response.Additional, TypeA)
	nsToCache := map[DnsCacheKey][]DnsRecord{}
	aToCache := map[DnsCacheKey][]DnsRecord{}
	ipv4s := []string{}
	for _, record := range nameserverRecords {
		cacheKey := DnsCacheKey{
			Qname:  record.Name,
			Qtype:  record.Type,
			Qclass: record.Class,
		}
		resolvedRecords, exists := nsToCache[cacheKey]
		if exists {
			resolvedRecords = append(resolvedRecords, record)
			nsToCache[cacheKey] = resolvedRecords
		} else {
			nsToCache[cacheKey] = []DnsRecord{record}
		}
	}
	for _, record := range nameserverIpv4Records {
		cacheKey := DnsCacheKey{
			Qname:  record.Name,
			Qtype:  record.Type,
			Qclass: record.Class,
		}
		resolvedRecords, exists := aToCache[cacheKey]
		if exists {
			resolvedRecords = append(resolvedRecords, record)
			aToCache[cacheKey] = resolvedRecords
		} else {
			aToCache[cacheKey] = []DnsRecord{record}
		}
		ipv4s = append(ipv4s, record.RData)
	}

	// cache the referral nameservers and their ip addresses
	for cacheKey, records := range nsToCache {
		r.cache.Put(cacheKey, DnsCacheValue{
			Answers:    records,
			Authority:  []DnsRecord{},
			Additional: []DnsRecord{},
		})
	}

	for cacheKey, records := range aToCache {
		r.cache.Put(cacheKey, DnsCacheValue{
			Answers:    records,
			Authority:  []DnsRecord{},
			Additional: []DnsRecord{},
		})
	}

	return ipv4s
}

func (r *RecursiveResolver) getCachedAuthoritativeNameserverIpv4s(domain string) ([]string, bool) {
	if domain == "." {
		return ROOT_NAME_SERVER_IPV4S, true
	}

	if !r.useCache() {
		return nil, false
	}

	// get nameservers
	nsCacheKey := DnsCacheKey{
		Qname:  domain,
		Qtype:  TypeNS,
		Qclass: ClassIN,
	}

	nsCacheEntry := r.cache.Get(nsCacheKey)
	if nsCacheEntry == nil {
		return nil, false
	}

	nameservers := []string{}
	for _, record := range nsCacheEntry.Value.Answers {
		nameservers = append(nameservers, record.RData)
	}

	// get the ip address for each nameserver
	ipv4s := []string{}
	for _, nameserver := range nameservers {
		aCacheKey := DnsCacheKey{
			Qname:  nameserver,
			Qtype:  TypeA,
			Qclass: ClassIN,
		}
		aCacheEntry := r.cache.Get(aCacheKey)
		if aCacheEntry == nil {
			continue
		}
		for _, record := range aCacheEntry.Value.Answers {
			ipv4s = append(ipv4s, record.RData)
		}
	}

	if len(ipv4s) == 0 {
		return nil, false
	}

	return ipv4s, true
}

func (r *RecursiveResolver) useCache() bool {
	return r.cache != nil
}

func (r *RecursiveResolver) getBestStartingNameserverIpv4s(qname string) []string {
	domains := getAllDomains(qname)
	fmt.Println("domains: ", domains)
	for _, domain := range domains {
		nameserverIpv4s, cached := r.getCachedAuthoritativeNameserverIpv4s(domain)
		if !cached {
			continue
		}
		return nameserverIpv4s
	}

	// will never reach here
	return []string{}
}

func (r *RecursiveResolver) cacheResponse(query DnsMessage, response DnsMessage) {
	if !r.useCache() {
		return
	}

	cacheKey := newDnsCacheKeyFromMessage(query)
	cacheValue := DnsCacheValue{
		Answers:    response.Answers,
		Authority:  response.Authority,
		Additional: response.Additional,
	}
	r.cache.Put(cacheKey, cacheValue)
}

func (r *RecursiveResolver) getCachedResponse(query DnsMessage) (*DnsMessage, bool) {
	if !r.useCache() {
		return nil, false
	}

	// get cached dns entry
	cacheKey := newDnsCacheKeyFromMessage(query)
	cacheEntry := r.cache.Get(cacheKey)
	if cacheEntry == nil {
		return nil, false
	}

	// build response dns message
	res := NewDnsMessageBuilder().
		From(query).
		Response().
		RecursionAvailable(true).
		Answers(cacheEntry.Value.Answers).
		Authority(cacheEntry.Value.Authority).
		Additional(cacheEntry.Value.Additional).
		Build()

	return &res, true
}

func getParentDomain(domain string) string {
	if domain == "." {
		return domain
	}
	parts := strings.Split(domain, ".")
	parents := parts[1:]
	parentDomain := strings.Join(parents, ".")
	if parentDomain == "" {
		parentDomain += "."
	}
	return parentDomain
}

func getAllDomains(domain string) []string {
	domains := []string{domain}
	for domain != "." {
		parent := getParentDomain(domain)
		domains = append(domains, parent)
		domain = parent
	}
	return domains
}

func filterRecordsByType(records []DnsRecord, qtype Type) []DnsRecord {
	filtered := []DnsRecord{}
	for _, record := range records {
		if record.Type == qtype {
			filtered = append(filtered, record)
		}
	}
	return filtered
}
