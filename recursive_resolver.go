package godns

import (
	"bytes"
	"fmt"
	"net"
)

type DNSRecursiveResolver struct {
	rootNameServer string
}

func NewDNSRecursiveResolver(rootNameServer string) *DNSRecursiveResolver {
	return &DNSRecursiveResolver{rootNameServer: rootNameServer}
}

func (r *DNSRecursiveResolver) Query(message *DNSMessage) (*DNSMessage, error) {
	fmt.Println("recursive resolve")

	res, err := r.send(r.rootNameServer, message)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *DNSRecursiveResolver) send(server string, message *DNSMessage) (*DNSMessage, error) {
	// connect to udp
	conn, err := net.Dial("udp", fmt.Sprintf("%v:%v", server, 53))
	if err != nil {
		return nil, err
	}

	// encode dns message to bytes
	var buffer bytes.Buffer
	encoder := &DNSMessageEncoder{message: message}
	encoder.Write(&buffer)

	// send udp datagram
	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		return nil, err
	}

	responseBuffer := make([]byte, 65535)
	n, err := conn.Read(responseBuffer)
	if err != nil {
		return nil, err
	}

	// parse dns message from buffer
	responseBuffer = responseBuffer[:n]
	parser := NewDNSMessageParser(responseBuffer)
	dnsMessage := DNSMessage{}
	err = parser.Parse(&dnsMessage)
	if err != nil {
		return nil, err
	}

	return &dnsMessage, nil
}
