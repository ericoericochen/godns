package godns

import (
	"fmt"
)

type RecursiveResolver struct {
	rootNameServer string
}

func NewRecursiveResolver(rootNameServer string) *RecursiveResolver {
	return &RecursiveResolver{rootNameServer: rootNameServer}
}

func (r *RecursiveResolver) Query(message *DnsMessage) (*DnsMessage, error) {
	res, err := r.resolve(message)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *RecursiveResolver) Resolve(name string, qtype Type, qclass Class) {

}

func (r *RecursiveResolver) resolve(message *DnsMessage) (*DnsMessage, error) {
	var err error
	var res *DnsMessage

	serverIP := r.rootNameServer

	for {
		fmt.Println(fmt.Sprintf("resolving query against nameserver: %v", serverIP))
		res, err = sendDnsMessage(serverIP, UDP_PORT, message)
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
