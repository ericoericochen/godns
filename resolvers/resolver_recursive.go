package resolvers

// import (
// 	"fmt"
// )

// type RecursiveResolver struct {
// 	rootNameServer string
// 	cache          Cache[DnsCacheKey, DnsCacheValue]
// }

// func NewRecursiveResolver(rootNameServer string) *RecursiveResolver {
// 	return &RecursiveResolver{rootNameServer: rootNameServer}
// }

// func (r *RecursiveResolver) Query(message *DnsMessage) (*DnsMessage, error) {
// 	res, err := r.resolve(message)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res, nil
// }

// func (r *RecursiveResolver) Resolve(qname string, qtype Type, qclass Class) (*DnsMessage, error) {
// 	// build dns query message
// 	message := NewDnsMessage()
// 	message.Header.QR = false
// 	message.Header.Opcode = OpcodeQUERY
// 	message.Header.Qdcount = 1
// 	message.Questions = append(message.Questions, &DNSQuestion{
// 		Qname:  qname,
// 		Qtype:  qtype,
// 		Qclass: qclass,
// 	})

// 	// recursively resolve dns query and cache results
// 	res, err := r.resolve(message)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return res, nil
// }

// func (r *RecursiveResolver) resolve(message *DnsMessage) (*DnsMessage, error) {
// 	var err error
// 	var res *DnsMessage

// 	serverIP := r.rootNameServer

// 	for {
// 		fmt.Println(fmt.Sprintf("resolving query against nameserver: %v", serverIP))
// 		res, err = sendDnsMessage(serverIP, UDP_PORT, message)
// 		if err != nil {
// 			return nil, err
// 		}
// 		fmt.Println("received response")
// 		fmt.Println(res)

// 		// did we receive the answer for the query
// 		if len(res.Answers) > 0 {
// 			fmt.Println("received answer")

// 			// cache the response

// 			break
// 		}

// 		fmt.Println("did not receive answer")

// 		var referralNS string
// 		for _, record := range res.Authority {
// 			if record.Type == TypeNS {
// 				referralNS = record.RData
// 				break
// 			}
// 		}

// 		fmt.Println("referral nameserver: ", referralNS)

// 		// find the ip address of the referral ns
// 		for _, record := range res.Additional {
// 			if record.Name == referralNS && record.Type == TypeA {
// 				fmt.Println("found ip address of ", referralNS, "=", record.RData)
// 				serverIP = record.RData
// 				break
// 			}
// 		}

// 		fmt.Println("----")
// 	}

// 	// cache results
// 	// r.cache.Put(DnsCacheKey{Qname: })

// 	return res, err
// }
