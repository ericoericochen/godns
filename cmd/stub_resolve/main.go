package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/ericoericochen/godns"
)

func main() {
	resolver := godns.NewStubResolver().WithCache(godns.NewDnsCache(1000))
	qname := "www.google.com"
	qtype := godns.TypeA
	qclass := godns.ClassIN

	// build query dns message
	query := godns.NewDnsQuery(qname, qtype, qclass, true)

	type resolveResult struct {
		response godns.DnsMessage
		err      error
		latency  time.Duration
	}

	results := make([]resolveResult, 3)
	for i := range results {
		start := time.Now()
		results[i].response = resolver.Resolve(query)
		results[i].latency = time.Since(start)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, "REQUEST\tLATENCY\tSTATUS")
	for i, result := range results {
		status := "ok"
		if result.err != nil {
			status = result.err.Error()
		}
		fmt.Fprintf(w, "%d\t%s\t%s\n", i+1, result.latency, status)
	}
	w.Flush()

	for i, result := range results {
		fmt.Printf("\nResponse %d:\n", i+1)
		if result.err != nil {
			fmt.Printf("  error: %v\n", result.err)
			continue
		}
		fmt.Println(result.response)
	}
}
