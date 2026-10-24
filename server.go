package godns

import (
	"bytes"
	"fmt"
	"log"
	"net"
)

type DnsServer struct {
	resolver DnsResolver
	conn     *net.UDPConn
}

func NewDnsServer(resolver DnsResolver) *DnsServer {
	server := &DnsServer{}
	server.resolver = resolver
	return server
}

func (s *DnsServer) Listen(port int) error {
	addr := &net.UDPAddr{Port: port}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	s.conn = conn

	log.Println(fmt.Sprintf("Dns Server listening on port %v", conn.LocalAddr()))

	buf := make([]byte, 512)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println("read error", err)
			continue
		}

		datagram := buf[:n]
		s.process(clientAddr, datagram)
	}
}

func (s *DnsServer) process(clientAddr *net.UDPAddr, datagram []byte) {
	fmt.Println(clientAddr)
	fmt.Println(datagram)

	// parse dns message
	var query DnsMessage
	decoder := NewDnsMessageDecoder(datagram)
	decoder.Decode(&query)

	// query resolver
	response := s.resolver.Resolve(query)

	fmt.Println(response)

	// encode response dns message
	var resBuffer bytes.Buffer
	encoder := NewDnsMessageEncoder(response)
	encoder.Write(&resBuffer)

	// send response datagram back to client
	_, err := s.conn.WriteToUDP(resBuffer.Bytes(), clientAddr)
	if err != nil {
		log.Println("write error: ", err)
	}

}
