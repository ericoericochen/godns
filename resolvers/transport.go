package resolvers

import (
	"bytes"
	"fmt"
	"net"

	"github.com/ericoericochen/godns"
)

const UDP_DNS_PORT = 53

func SendUDPDatagram(serverIp string, port int, datagram []byte) ([]byte, error) {
	// create udp connection
	conn, err := net.Dial("udp", fmt.Sprintf("%v:%v", serverIp, port))
	if err != nil {
		return nil, err
	}

	// send datagram
	_, err = conn.Write(datagram)

	// receive response
	resBuffer := make([]byte, 512)
	n, err := conn.Read(resBuffer)
	if err != nil {
		return nil, err
	}

	return resBuffer[:n], nil
}

// Send the DnsMessage to UDP nameserver
// The datagram is a valid encoded DnsMessage
func SendDnsMessage(nameserverIp string, port int, message godns.DnsMessage) godns.DnsMessage {
	// convert dns message to bytes
	var datagram bytes.Buffer
	encoder := godns.NewDnsMessageEncoder(message)
	encoder.Write(&datagram)

	// send the dns message udp datagram
	resBuffer, err := SendUDPDatagram(nameserverIp, port, datagram.Bytes())

	// create message builder
	builder := godns.NewDnsMessageBuilder().From(message)
	if err != nil {
		return builder.Rcode(godns.RcodeServerFailure).Build()
	}

	// parse dns message
	res := godns.DnsMessage{}
	decoder := godns.NewDnsMessageDecoder(resBuffer)
	err = decoder.Decode(&res)
	if err != nil {
		return builder.Rcode(godns.RcodeServerFailure).Build()
	}

	return res
}
