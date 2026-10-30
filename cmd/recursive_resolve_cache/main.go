package main

import (
	"fmt"
	"time"

	"github.com/ericoericochen/godns"
)

func main() {
	resolver := godns.NewRecursiveResolver().WithRootNameserver(godns.E_ROOT_NAMESERVER).WithCache(godns.NewDnsCache(1000))

	resolve := func(qname string) {
		query := godns.NewDnsQuery(qname, godns.TypeA, godns.ClassIN, false)

		start := time.Now()
		response := resolver.Resolve(query)
		latency := time.Since(start)

		fmt.Printf("%-16s latency: %-12s rcode: %v answers: %d\n",
			qname, latency, response.Header.Rcode, len(response.Answers))
		fmt.Println(response)
	}

	resolve("www.google.com.")
	resolve("www.google.com.")
	resolve("google.com.")
	resolve("cloud.google.com.")
	resolve("news.google.com.")
}
