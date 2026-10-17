package godns

import (
	"bytes"
	"fmt"
	"net"
)

const UDP_PORT = 53

func sendDnsMessage(serverIP string, port int, message *DnsMessage) (*DnsMessage, error) {
	// connect to udp
	conn, err := net.Dial("udp", fmt.Sprintf("%v:%v", serverIP, port))
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

	// receive response
	responseBuffer := make([]byte, 512)
	n, err := conn.Read(responseBuffer)
	if err != nil {
		return nil, err
	}

	// parse dns message from buffer
	responseBuffer = responseBuffer[:n]
	parser := NewDNSMessageParser(responseBuffer)
	response := DnsMessage{}
	err = parser.Parse(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
