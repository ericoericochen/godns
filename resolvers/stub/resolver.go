package stub

import (
	"bytes"

	"github.com/ericoericochen/godns"
)

type StubResolver struct {
	server string
}

func New(server string) godns.DnsResolver {
	resolver := &StubResolver{server: server}
	return resolver
}

func (r *StubResolver) Resolve(message *godns.DnsMessage) (*godns.DnsMessage, error) {
	// validate dns query message
	err := godns.ValidateQueryDnsMessage(message)
	if err != nil {
		return nil, err
	}

	// encode dns message to bytes
	var buffer bytes.Buffer
	encoder := godns.NewDnsMessageEncoder(message)
	encoder.Write(&buffer)

	// send the dns message
	res, err := godns.SendDnsMessage(r.server, godns.UDP_DNS_PORT, buffer.Bytes())
	return res, err
}
