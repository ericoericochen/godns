package main

import (
	"fmt"

	"github.com/ericoericochen/godns"
)

func main() {
	cache := godns.NewDnsCache(2)

	firstKey := godns.DnsCacheKey{Qname: "first.example", Qtype: godns.TypeA, Qclass: godns.ClassIN}
	secondKey := godns.DnsCacheKey{Qname: "second.example", Qtype: godns.TypeA, Qclass: godns.ClassIN}
	thirdKey := godns.DnsCacheKey{Qname: "third.example", Qtype: godns.TypeA, Qclass: godns.ClassIN}

	firstValue := godns.DnsCacheValue{
		TTL: 100,
		Answers: []godns.DnsRecord{
			{Name: "first.example", Type: godns.TypeA, Class: godns.ClassIN, TTL: 100, RData: "192.0.2.1"},
		},
	}
	secondValue := godns.DnsCacheValue{
		TTL: 200,
		Answers: []godns.DnsRecord{
			{Name: "second.example", Type: godns.TypeA, Class: godns.ClassIN, TTL: 200, RData: "192.0.2.2"},
		},
	}
	thirdValue := godns.DnsCacheValue{
		TTL: 300,
		Answers: []godns.DnsRecord{
			{Name: "third.example", Type: godns.TypeA, Class: godns.ClassIN, TTL: 300, RData: "192.0.2.3"},
		},
	}

	cache.Put(firstKey, firstValue)
	cache.Put(secondKey, secondValue)

	entry := cache.Get(firstKey)
	fmt.Printf("first access to first key: %+v\n", entry)

	entry = cache.Get(firstKey)
	fmt.Printf("second access to first key: %+v\n", entry)

	cache.Put(thirdKey, thirdValue)

	firstEntry := cache.Get(firstKey)
	secondEntry := cache.Get(secondKey)
	thirdEntry := cache.Get(thirdKey)

	fmt.Printf("first key after third insert:  %+v\n", firstEntry)
	fmt.Printf("second key after third insert: %+v\n", secondEntry)
	fmt.Printf("third key after third insert:  %+v\n", thirdEntry)
}
