package godns

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
	return r.resolve(query)
}

func (r *RecursiveResolver) resolve(query DnsMessage) DnsMessage {
	// validate query
	question := query.Questions[0]
	name := question.Qname

	// get cached response

	return r.recursiveResolve(name, query)
}

func (r *RecursiveResolver) recursiveResolve(domain string, query DnsMessage) DnsMessage {
	// domain - owner of the nameserver responsible for the domain

	// get the authoritative nameservers for domain
	nameservers := getAuthoritativeNameservers(domain)
	if nameservers == nil {
		// resolve the query against the nameserver for the parent domain
		parentDomain := getParentDomain(domain)
		res := r.recursiveResolve(parentDomain, query)
		if len(res.Answers) > 0 {
			return res
		} else {
			// FIX
			nameservers := getNSRdatas(res.Authority)
			nameserver := nameservers[0]
			return SendDnsMessage(nameserver, UDP_DNS_PORT, query)
		}

	} else {
		nameserver := nameservers[0]
		return SendDnsMessage(nameserver, UDP_DNS_PORT, query)
	}
}

func getParentDomain(domain string) string {
	return ""
}

// get the ipv4 address for the authoritative nameserver for domain
func getAuthoritativeNameservers(domain string) []string {
	if domain == "." {
		return ROOT_NAME_SERVERS
	}

	// get the nameservers for the domain

	// get ipv4 for the nameservers

	return []string{}
}

// func (r *RecursiveResolver) resolve(query DnsMessage) DnsMessage {
// 	serverIp := rootNameServer
// 	for {
// 		fmt.Println(fmt.Sprintf("resolving query against nameserver: %v", rootNameServer))
// 		res := SendDnsMessage(serverIp, UDP_DNS_PORT, query)
// 		fmt.Println(res)
// 		if res.Header.Rcode != RcodeNoError {
// 			return res
// 		}

// 		// cache response
// 		valuesToCache := map[DnsCacheKey]*DnsCacheValue{}
// 		records := []DnsRecord{}
// 		records = append(records, res.Answers...)
// 		records = append(records, res.Authority...)
// 		records = append(records, res.Additional...)

// 		for _, record := range records {
// 			cacheKey := DnsCacheKey{
// 				Qname:  record.Name,
// 				Qtype:  record.Type,
// 				Qclass: record.Class,
// 			}
// 			cacheValue, exists := valuesToCache[cacheKey]
// 			if exists {
// 				cacheValue.Answers = append(cacheValue.Answers, record)
// 			} else {
// 				cacheValue := NewDnsCacheValue()
// 				cacheValue.Answers = append(cacheValue.Answers, record)
// 				valuesToCache[cacheKey] = &cacheValue
// 			}
// 		}

// 		fmt.Println("values to cache")
// 		fmt.Println(valuesToCache)

// 		// did we receive answer to query vs. referral
// 		if len(res.Answers) > 0 {
// 			fmt.Println("received answer")
// 			return res
// 		}

// 		fmt.Println("did not receive answer")
// 		referralNameservers := getNSRdatas(res.Authority)
// 		referralNameserver := referralNameservers[0]

// 		fmt.Println("referral nameserver: ", referralNameservers[0])

// 		// find the ip address of the referral ns
// 		for _, record := range res.Additional {
// 			if record.Name == referralNameserver && record.Type == TypeA {
// 				fmt.Println("found ip address of ", referralNameserver, "=", record.RData)
// 				serverIp = record.RData
// 				break
// 			}
// 		}
// 	}
// }

func getNSRdatas(records []DnsRecord) []string {
	rdatas := []string{}
	for _, record := range records {
		if record.Type == TypeNS {
			rdatas = append(rdatas, record.RData)
		}
	}
	return rdatas
}
