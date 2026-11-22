package main

import (
	"errors"
	"fmt"
	"log"
	"net/netip"
	"os"
	"strconv"

	"github.com/ericoericochen/godns"
)

const SERVER_CONFIG_FILE = "./server.config.json"
const DNS_RESOLVER_TYPE_ENV = "DNS_RESOLVER_TYPE"
const DNS_CACHE_CAPACITY_ENV = "DNS_CACHE_CAPACITY"
const DNS_ROOT_NAMESERVER_ENV = "DNS_ROOT_NAMESERVER"
const DNS_STUB_NAMESERVER_ENV = "DNS_STUB_NAMESERVER"

const DEFAULT_DNS_ROOT_NAMESERVER = "a.root-servers.net."

type resolverType string

const (
	stub      resolverType = "stub"
	recursive resolverType = "recursive"
)

var invalidDnsResolverType = errors.New("invalid dns resolver type")

func getDNSResolverType() resolverType {
	value, exists := os.LookupEnv(DNS_RESOLVER_TYPE_ENV)
	if !exists {
		return recursive
	}
	return resolverType(value)
}

func getDNSCacheCapacity() int {
	value, exists := os.LookupEnv(DNS_CACHE_CAPACITY_ENV)
	if !exists {
		return 1000
	}
	capacity, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return capacity
}

func getDNSRootNameserver() string {
	value, exists := os.LookupEnv(DNS_ROOT_NAMESERVER_ENV)
	if !exists {
		return DEFAULT_DNS_ROOT_NAMESERVER
	}
	return value
}

func getDNSStubNameserver() string {
	value, exists := os.LookupEnv(DNS_STUB_NAMESERVER_ENV)
	if !exists {
		return godns.CLOUDFLARE_DNS_SERVER
	}
	return value
}

func buildDnsResolverFromEnv() (godns.DnsResolver, error) {
	rt := getDNSResolverType()
	cacheCapacity := getDNSCacheCapacity()

	switch rt {
	case stub:
		{
			server := getDNSStubNameserver()
			serverIp, err := netip.ParseAddr(server)
			if err != nil {
				return nil, err
			}
			log.Println(fmt.Sprintf("initializing stub resolver with server=%v, cacheCapacity=%v", serverIp, cacheCapacity))
			resolver := godns.NewStubResolver().
				WithServer(serverIp).
				WithCache(godns.NewDnsCache(cacheCapacity))
			return resolver, nil
		}
	case recursive:
		{
			rootNameserver := getDNSRootNameserver()
			resolver := godns.NewRecursiveResolver().
				WithRootNameserver(rootNameserver).
				WithCache(godns.NewDnsCache(cacheCapacity))
			return resolver, nil
		}
	}

	return nil, invalidDnsResolverType
}

func buildDnsServerFromEnv() (*godns.DnsServer, error) {
	resolver, err := buildDnsResolverFromEnv()
	if err != nil {
		log.Println("Failed to initialize dns resolver")
		return nil, err
	}
	server := godns.NewDnsServer(resolver)
	return server, nil
}

func main() {
	// build the server
	server, err := buildDnsServerFromEnv()

	if err != nil {
		log.Println("Failed to initialize dns server")
		return
	}

	// start the server on port 53
	server.Listen(godns.UDP_DNS_PORT)
}
