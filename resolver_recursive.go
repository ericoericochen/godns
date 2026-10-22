package godns

import (
	"fmt"
	"strings"
)

type RecursiveResolver struct {
	cache *DnsCache
}

func NewRecursiveResolver() *RecursiveResolver {
	return &RecursiveResolver{}
}

const rootNameServer = "198.41.0.4"

func (r *RecursiveResolver) WithCache(cache *DnsCache) *RecursiveResolver {
	r.cache = cache
	return r
}

func (r *RecursiveResolver) Resolve(query DnsMessage) DnsMessage {
	// validate dns message
	return r.resolve(query)
}

func (r *RecursiveResolver) resolve(query DnsMessage) DnsMessage {
	fmt.Println("yo")
	// find the best starting nameserver to resolve against
	nameservers := r.getBestStartingNameserverIpv4s(query.Qname())

	fmt.Println(nameservers)

	serverIp := rootNameServer
	for {
		fmt.Println(fmt.Sprintf("resolving query against nameserver: %v", rootNameServer))
		res := SendDnsMessage(serverIp, UDP_DNS_PORT, query)
		fmt.Println(res)
		if res.Header.Rcode != RcodeNoError {
			return res
		}

		// did we receive answer to query vs. referral
		if len(res.Answers) > 0 {
			fmt.Println("received answer")
			return res
		}

		fmt.Println("did not receive answer")
		referralNameservers := getNSRdatas(res.Authority)
		referralNameserver := referralNameservers[0]

		fmt.Println("referral nameserver: ", referralNameservers[0])

		// find the ip address of the referral ns
		for _, record := range res.Additional {
			if record.Name == referralNameserver && record.Type == TypeA {
				fmt.Println("found ip address of ", referralNameserver, "=", record.RData)
				serverIp = record.RData
				break
			}
		}
	}
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
	fmt.Println(domain)
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

func getNSRdatas(records []DnsRecord) []string {
	rdatas := []string{}
	for _, record := range records {
		if record.Type == TypeNS {
			rdatas = append(rdatas, record.RData)
		}
	}
	return rdatas
}
