package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/netip"
	"os"

	"github.com/ericoericochen/godns"
)

const SERVER_CONFIG_FILE = "./server.config.json"
const RESOLVER_TYPE_ENV = "RESOLVER_TYPE"

type resolverType string

const (
	stub      resolverType = "stub"
	recursive resolverType = "recursive"
)

func getResolverType() resolverType {
	value, exists := os.LookupEnv(RESOLVER_TYPE_ENV)
	if !exists {
		return recursive
	}
	return resolverType(value)
}

func buildDnsResolver(config map[string]any, rt resolverType) (godns.DnsResolver, error) {
	resolverConfig := config[string(rt)].(map[string]any)

	switch rt {
	case stub:
		{
			server := resolverConfig["server"].(string)
			cacheCapacity := int(resolverConfig["cache"].(map[string]any)["capacity"].(float64))
			serverIp, err := netip.ParseAddr(server)
			if err != nil {
				return nil, err
			}
			resolver := godns.NewStubResolver().
				WithServer(serverIp).
				WithCache(godns.NewDnsCache(cacheCapacity))
			return resolver, nil
		}
	case recursive:
		{

		}
	}

	return nil, nil
}

func buildDnsServer(config map[string]any, rt resolverType) (*godns.DnsServer, error) {
	resolver, err := buildDnsResolver(config, rt)
	if err != nil {
		log.Println("Failed to initialize dns resolver")
		return nil, err
	}
	server := godns.NewDnsServer(resolver)
	return server, nil
}

func main() {
	fmt.Println(SERVER_CONFIG_FILE)

	// parse the server config file
	f, err := os.Open(SERVER_CONFIG_FILE)
	if err != nil {
		log.Println("Failed to start server. Unable to open server.config.json.")
	}
	defer f.Close()
	var config map[string]any

	if err := json.NewDecoder(f).Decode(&config); err != nil {
		log.Println("Failed to start server. Unable to parse server.config.json.")
	}

	// get the resolver type
	rt := getResolverType()

	// build the server
	server, err := buildDnsServer(config, rt)

	if err != nil {
		log.Println("Failed to initialize dns server")
		return
	}

	server.Listen(godns.UDP_DNS_PORT)

	// resolver := godns.NewStubResolver()
	// server := godns.NewDnsServer(resolver)
	// server.Listen(godns.UDP_DNS_PORT)
}
