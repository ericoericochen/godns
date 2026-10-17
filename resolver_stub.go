package godns

import (
	"bytes"
	"fmt"
	"net"
)

type StubResolver struct {
	server string
}

func NewStubResolver(server string) *StubResolver {
	return &StubResolver{server: server}
}

const port = 53

func (r *StubResolver) Query(message *DnsMessage) (*DnsMessage, error) {
	// connect to udp dns server
	conn, err := net.Dial("udp", fmt.Sprintf("%v:%v", r.server, port))
	if err != nil {
		return nil, err
	}

	// encode dns message
	var buffer bytes.Buffer

	encoder := &DNSMessageEncoder{message: message}
	encoder.Write(&buffer)

	// send dns message
	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		return nil, err
	}

	// receive the response in a buffer
	responseBuffer := make([]byte, 65535)
	n, err := conn.Read(responseBuffer)
	if err != nil {
		return nil, err
	}

	// parse dns message from buffer
	responseBuffer = responseBuffer[:n]
	parser := NewDNSMessageParser(responseBuffer)
	dnsMessage := DnsMessage{}
	err = parser.Parse(&dnsMessage)
	if err != nil {
		return nil, err
	}

	return &dnsMessage, nil
}
