package main

import (
	"bytes"
	"fmt"

	"github.com/ericoericochen/godns"
)

func main() {
	qname := "www.google.com."
	qtype := godns.TypeA
	qclass := godns.ClassIN
	message := godns.NewDnsQuery(qname, qtype, qclass, true)
	message.Header.Ancount = 1
	message.Answers = []godns.DnsRecord{
		godns.DnsRecord{
			Name:  qname,
			Type:  qtype,
			Class: qclass,
			TTL:   300,
			RData: "1.1.1.1",
		},
	}

	var buf bytes.Buffer
	encoder := godns.NewDnsMessageEncoder(message)
	encoder.Write(&buf)
	fmt.Println(buf.Bytes())
	fmt.Println("length: ", len(buf.Bytes()))

	var decoded godns.DnsMessage
	decoder := godns.NewDnsMessageDecoder(buf.Bytes())
	err := decoder.Decode(&decoded)

	fmt.Println("err: ", err)

	fmt.Println("decoded: ")
	fmt.Println(decoded)
}
