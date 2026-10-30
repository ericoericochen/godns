package main

import (
	"fmt"

	"github.com/ericoericochen/godns"
)

func main() {
	resolver := godns.NewRecursiveResolver().WithRootNameserver(godns.E_ROOT_NAMESERVER)
	qname := "en.wikipedia.org."
	// qname := "www.google.com."
	qtype := godns.TypeA
	qclass := godns.ClassIN

	// build query dns message
	message := godns.NewDnsQuery(qname, qtype, qclass, false)
	res := resolver.Resolve(message)
	fmt.Println("------")

	fmt.Println(res)
}
