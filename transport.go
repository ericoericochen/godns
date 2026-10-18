package godns

import (
	"fmt"
	"net"
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
func SendDnsMessage(nameserverIp string, port int, datagram []byte) (*DnsMessage, error) {
	resBuffer, err := SendUDPDatagram(nameserverIp, port, datagram)

	// parse dns message
	decoder := NewDnsMessageDecoder(resBuffer)
	res := DnsMessage{}
	err = decoder.Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
