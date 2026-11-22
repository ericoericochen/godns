package godns

import (
	"cmp"
	"fmt"
	"log"
	"net/netip"
	"slices"
)

type nameserverData struct {
	name  string
	ipv4s []netip.Addr
}

type RecursiveResolver struct {
	cache          *DnsCache
	rootNameserver string
}

const defaultRootNamserver = A_ROOT_NAMESERVER

func NewRecursiveResolver() *RecursiveResolver {
	return &RecursiveResolver{
		rootNameserver: defaultRootNamserver,
	}
}

func (r *RecursiveResolver) WithRootNameserver(rootNameserver string) *RecursiveResolver {
	r.rootNameserver = rootNameserver
	return r
}

func (r *RecursiveResolver) WithCache(cache *DnsCache) *RecursiveResolver {
	r.cache = cache
	return r
}

func (r *RecursiveResolver) Resolve(query DnsMessage) DnsMessage {
	// validate dns query message
	err := validateQueryDnsMessage(query)
	if err != nil {
		return toResponse(query, RcodeFormatError)
	}

	// validate dns message
	return r.resolve(query)
}

// Implements the DNS resolver algorithm
// https://www.rfc-editor.org/info/rfc1034/#section-5.3.3
func (r *RecursiveResolver) resolve(query DnsMessage) DnsMessage {
	// check if query response is cached
	res, cached := r.getCachedResponse(query)
	if cached {
		log.Println("cache hit")
		return *res
	}

	// find the best starting nameserver to resolve against
	nameservers := r.getBestStartingNameservers(query.Qname())

	// sort the nameservers to begin with nameserver with highest count of ipv4 addresses
	slices.SortFunc(nameservers, func(a, b nameserverData) int {
		return cmp.Compare(len(b.ipv4s), len(a.ipv4s))
	})

	// iterate until we have received authoritative answer to dns query
	for {
		// query against nameserver sequentially, returning the first
		// non-error response
		res := r.queryNameservers(query, nameservers)
		if res.Header.Rcode != RcodeNoError {
			break
		}

		fmt.Println(res)

		// received answer
		if len(res.Answers) > 0 {
			cnameRecords := filterRecordsByType(res.Answers, TypeCNAME)

			// if there are no cname records, we have got the answer to the query
			if len(cnameRecords) == 0 {
				res.Header.RecursionAvailable = true
				return res
			}

			cname := cnameRecords[0].RData
			cnameQuery := NewDnsQuery(cname, query.Qtype(), query.Qclass())
			fmt.Println("cname query: ")
			fmt.Println(cnameQuery)
			cnameRes := r.resolve(cnameQuery)

			fmt.Println("cname res: ")
			fmt.Println(cnameRes)
			builder := NewDnsMessageBuilder().From(res).RecursionAvailable(true)
			builder = builder.
				Answers(cnameRes.Answers)

			finalRes := builder.Build()
			finalRes.Header.Nscount = 0
			finalRes.Authority = []DnsRecord{}
			r.cacheResponse(query, finalRes)
			return finalRes
		}

		// get the next set of nameservers to query
		nextNameservers := r.getNextNameservers(res)
		nameservers = nextNameservers
		fmt.Println("next nameservers: ")
		fmt.Println(nameservers)
	}

	return toServerFailureResponse(query)
}

func (r *RecursiveResolver) queryNameservers(query DnsMessage, nameservers []nameserverData) DnsMessage {
	for _, nameserver := range nameservers {
		numIpv4s := len(nameserver.ipv4s)
		if numIpv4s == 0 {
			fmt.Println("no ipv4s for nameserver")
			nameserverQuery := NewDnsQuery(nameserver.name, TypeA, ClassIN)
			nameserverRes := r.resolve(nameserverQuery)
			fmt.Println("nameserver res: ")
			fmt.Println(nameserverRes)

			nameserverIps := getIpv4sFromRecords(nameserverRes.Answers)
			nameserver.ipv4s = nameserverIps
			fmt.Println("nameserver ips: ", nameserverIps)
		}
		for _, nsIp := range nameserver.ipv4s {
			log.Println(fmt.Sprintf("resolving qname=%v qtype=%v qclass=%v with ns %v %v",
				query.Qname(),
				typeNameMap[query.Qtype()],
				classNameMap[query.Qclass()],
				nameserver.name,
				nsIp))

			// query the dns nameserver
			res := SendDnsMessage(nsIp, UDP_DNS_PORT, query)

			// if fail, try the next nameserver ip
			if res.Header.Rcode != RcodeNoError {
				log.Println(fmt.Sprintf("error from resolving qname=%v qtype=%v qclass=%v with ns %v %v",
					query.Qname(),
					typeNameMap[query.Qtype()],
					classNameMap[query.Qclass()],
					nameserver.name,
					nsIp))
				continue
			}

			return res
		}
	}
	return toServerFailureResponse(query)
}

func (r *RecursiveResolver) getBestStartingNameservers(qname string) []nameserverData {
	// start at qname, then parent, grandparent towards root
	// to find the best starting nameserver to resolve against
	domains := getAllDomains(qname)
	for _, domain := range domains {
		nameservers, exists := r.getNameservers(domain)
		if !exists {
			continue
		}
		return nameservers
	}

	// will never reach here
	return []nameserverData{}
}

func (r *RecursiveResolver) getNameservers(domain string) ([]nameserverData, bool) {
	if domain == "." {
		nameserver := nameserverData{
			name:  r.rootNameserver,
			ipv4s: []netip.Addr{rootNameserverIpMap[r.rootNameserver]},
		}
		return []nameserverData{nameserver}, true
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

	nameservers := []nameserverData{}

	for _, record := range nsCacheEntry.Value.Answers {
		name := record.RData
		ipv4s := []netip.Addr{}

		// get cached A records for the nameserver
		aCacheKey := DnsCacheKey{
			Qname:  name,
			Qtype:  TypeA,
			Qclass: ClassIN,
		}
		aCacheEntry := r.cache.Get(aCacheKey)
		if aCacheEntry != nil {
			for _, aRecord := range aCacheEntry.Value.Answers {
				ip, err := netip.ParseAddr(aRecord.RData)
				if err == nil {
					ipv4s = append(ipv4s, ip)
				}
			}
		}

		nameserver := nameserverData{name: name, ipv4s: ipv4s}
		nameservers = append(nameservers, nameserver)
	}

	if len(nameservers) == 0 {
		return nil, false
	}

	return nameservers, true
}

// // response doesn't have answer and refers to next nameservers to query
func (r *RecursiveResolver) getNextNameservers(response DnsMessage) []nameserverData {
	// get nameserver records and their A records
	nameserverRecords := filterRecordsByType(response.Authority, TypeNS)
	nameserverIpv4Records := filterRecordsByType(response.Additional, TypeA)
	nsToCache := map[DnsCacheKey][]DnsRecord{}
	aToCache := map[DnsCacheKey][]DnsRecord{}

	nsToNameservers := map[string]*nameserverData{}
	nameservers := []nameserverData{}

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
		nsToNameservers[record.RData] = &nameserverData{name: record.RData, ipv4s: []netip.Addr{}}
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
		ip, err := netip.ParseAddr(record.RData)
		if err == nil {
			nameserver, exists := nsToNameservers[record.Name]
			if exists {
				ipv4s := append(nameserver.ipv4s, ip)
				nameserver.ipv4s = ipv4s
			}
		}
	}

	// materialize the nameservers slice
	for _, record := range nameserverRecords {
		name := record.RData
		nameserver := nsToNameservers[name]
		nameservers = append(nameservers, *nameserver)
	}

	if !r.useCache() {
		return nameservers
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

	return nameservers
}

// cache methods

func (r *RecursiveResolver) useCache() bool {
	return r.cache != nil
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
