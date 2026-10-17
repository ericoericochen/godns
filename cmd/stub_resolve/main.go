package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/ericoericochen/godns"
)

func main() {
	server := "1.1.1.1"
	resolver := godns.NewStubResolver(server)
	qname := "www.google.com"
	qtype := godns.TypeA
	qclass := godns.ClassIN

	// build query dns message
	message := godns.NewDnsMessage()
	id := uint16(rand.UintN(1 << 16))
	// var id uint16 = 2
	message.Header.ID = id
	message.Header.QR = false
	message.Header.RecursionDesired = true
	message.Header.Qdcount = 1
	message.Questions = append(message.Questions, &godns.DNSQuestion{
		Qname:  qname,
		Qtype:  qtype,
		Qclass: qclass,
	})

	fmt.Println(message)

	res, err := resolver.Query(message)
	fmt.Println(res, err)
}
