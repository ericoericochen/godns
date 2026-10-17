package internal

import "github.com/ericoericochen/godns"

func SendUDPDatagram(serverIp string, port int, datagram []byte) ([]byte, error) {
	return nil, nil
}

// Send the DnsMessage to UDP nameserver
// The datagram is a valid encoded DnsMessage
func SendDnsMessage(nameserverIp string, port int, datagram []byte) (*godns.DnsMessage, error) {
	return nil, nil
}
