package main

import (
	"fmt"

	"github.com/ericoericochen/godns"
	"github.com/ericoericochen/godns/resolvers/stub"
)

func main() {
	server := "1.1.1.1"
	resolver := stub.New(server)
	qname := "www.google.com"
	qtype := godns.TypeA
	qclass := godns.ClassIN

	// build query dns message
	query := godns.NewDnsQuery(qname, qtype, qclass, true)

	fmt.Println(query)

	res, err := resolver.Resolve(query)
	fmt.Println(err)
	fmt.Println(res)
}
