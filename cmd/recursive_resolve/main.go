package main

import (
	"fmt"
	"net/netip"

	"github.com/ericoericochen/godns"
)

func main() {
	yo := netip.MustParseAddr(godns.E_ROOT_NAMESERVER_IPV4)

	fmt.Println("yo: ", yo)

	resolver := godns.NewRecursiveResolver().WithRootNameserver(godns.E_ROOT_NAMESERVER)
	qname := "www.google.com."
	qtype := godns.TypeA
	qclass := godns.ClassIN

	// build query dns message
	message := godns.NewDnsQuery(qname, qtype, qclass, false)
	res := resolver.Resolve(message)
	fmt.Println("------")

	fmt.Println(res.Header)
	fmt.Println(res.Questions[0])
	fmt.Println(res)
}
