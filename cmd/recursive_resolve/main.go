package main

import (
	"fmt"

	"github.com/ericoericochen/godns"
)

func main() {
	rootNameServer := "192.203.230.10"
	resolver := godns.NewDNSRecursiveResolver(rootNameServer)
	qname := "www.google.com"
	qtype := godns.TypeA
	qclass := godns.ClassIN

	// build query dns message
	message := godns.NewDNSMessage()
	// id := uint16(rand.UintN(1 << 16))
	var id uint16 = 2
	message.Header.ID = id
	message.Header.QR = false
	message.Header.RecursionDesired = true
	message.Header.Qdcount = 1
	message.Questions = append(message.Questions, &godns.DNSQuestion{
		Qname:  qname,
		Qtype:  qtype,
		Qclass: qclass,
	})

	res, err := resolver.Query(message)
	fmt.Println(res, err)
}
