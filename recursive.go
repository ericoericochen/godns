package godns

import "fmt"

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

func getNSRdatas(records []DnsRecord) []string {
	rdatas := []string{}
	for _, record := range records {
		if record.Type == TypeNS {
			rdatas = append(rdatas, record.RData)
		}
	}
	return rdatas
}
