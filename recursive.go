package godns

// import (
// 	"fmt"

// 	"github.com/ericoericochen/godns"
// )

// type RecursiveResolver struct {
// 	cache *godns.DnsCache
// }

// func NewRecursiveResolver() DnsResolver {
// 	return &RecursiveResolver{}
// }

// const rootNameServer = "198.41.0.4"

// func (r *RecursiveResolver) WithCache(cache *godns.DnsCache) DnsResolver {
// 	r.cache = cache
// 	return r
// }

// func (r *RecursiveResolver) Resolve(message godns.DnsMessage) godns.DnsMessage {
// 	return r.resolve(message)
// }

// func (r *RecursiveResolver) resolve(message godns.DnsMessage) godns.DnsMessage {
// 	serverIp := rootNameServer
// 	for {
// 		fmt.Println(fmt.Sprintf("resolving query against nameserver: %v", rootNameServer))
// 		res := SendDnsMessage(serverIp, UDP_DNS_PORT, message)
// 		fmt.Println(res)
// 		if res.Header.Rcode != godns.RcodeNoError {
// 			return res
// 		}

// 		// did we receive answer to query vs. referral
// 		if len(res.Answers) > 0 {
// 			fmt.Println("received answer")
// 			return res
// 		}

// 		fmt.Println("did not receive answer")
// 		var referralNS string
// 		for _, record := range res.Authority {
// 			if record.Type == godns.TypeNS {
// 				referralNS = record.RData
// 				break
// 			}
// 		}

// 		fmt.Println("referral nameserver: ", referralNS)

// 		// find the ip address of the referral ns
// 		for _, record := range res.Additional {
// 			if record.Name == referralNS && record.Type == godns.TypeA {
// 				fmt.Println("found ip address of ", referralNS, "=", record.RData)
// 				serverIp = record.RData
// 				break
// 			}
// 		}
// 	}
// }
