package main

import (
	"fmt"

	"github.com/ericoericochen/godns"
)

func main() {
	// rootNameServer := "192.203.230.10"
	resolver := godns.NewRecursiveResolver()
	qname := "google.com."
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
