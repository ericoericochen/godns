package main

import "github.com/ericoericochen/godns"

func main() {
	resolver := godns.NewStubResolver()
	server := godns.NewDnsServer(resolver)
	server.Listen(53)
}
