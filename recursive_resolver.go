package godns

import (
	"fmt"
)

type DNSRecursiveResolver struct {
	rootNameServer string
}

func NewDNSRecursiveResolver(rootNameServer string) *DNSRecursiveResolver {
	return &DNSRecursiveResolver{rootNameServer: rootNameServer}
}

func (r *DNSRecursiveResolver) Query(message *DNSMessage) (*DNSMessage, error) {
	res, err := r.resolve(message)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *DNSRecursiveResolver) resolve(message *DNSMessage) (*DNSMessage, error) {
	var err error
	var res *DNSMessage

	serverIP := r.rootNameServer

	for {
		fmt.Println(fmt.Sprintf("resolving query against nameserver: %v", serverIP))
		res, err = sendDNSMessage(serverIP, UDP_PORT, message)
		if err != nil {
			return nil, err
		}
		fmt.Println("received response")
		fmt.Println(res)

		// did we receive the answer for the query
		if len(res.Answers) > 0 {
			fmt.Println("received answer")
			break
		}

		fmt.Println("did not receive answer")

		var referralNS string
		for _, record := range res.Authorities {
			if record.Type == TypeNS {
				referralNS = record.RData
				break
			}
		}

		fmt.Println("referral nameserver: ", referralNS)

		// find the ip address of the referral ns
		for _, record := range res.Additionals {
			if record.Name == referralNS && record.Type == TypeA {
				fmt.Println("found ip address of ", referralNS, "=", record.RData)
				serverIP = record.RData
				break
			}
		}

		fmt.Println("----")
	}

	return res, err
}
