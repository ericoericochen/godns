package resolvers

import "github.com/ericoericochen/godns"

type DnsResolver interface {
	Resolve(*godns.DnsMessage) (*godns.DnsMessage, error)
}
