package main

import "github.com/ericoericochen/godns"

func main() {
	// resolver := godns.NewStubResolver()
	resolver := godns.NewRecursiveResolver()
	server := godns.NewDnsServer(resolver)
	server.Listen(53)
}
