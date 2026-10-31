package godns

import (
	"fmt"
	"log"
	"net"
	"time"
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

	log.Println(fmt.Sprintf("dns server listening on port %v", conn.LocalAddr()))

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
	log.Println(fmt.Sprintf("received dns query from %v", clientAddr))

	// parse dns message
	start := time.Now()
	query, err := Decode(datagram)
	if err != nil {
		log.Println("error decoding dns query")
		return
	}

	log.Println(fmt.Sprintf("qname=%v, qtype=%v, qclass=%v",
		query.Qname(),
		typeNameMap[query.Qtype()],
		classNameMap[query.Qclass()]))

	// query resolver
	response := s.resolver.Resolve(query)
	elapsed := time.Since(start)

	log.Println(fmt.Sprintf("took %v", elapsed))
	log.Println(fmt.Sprintf("response rcode=%v, answers=%v, authority=%v, additional=%v",
		rcodeMap[response.Header.Rcode],
		len(response.Answers),
		len(response.Authority),
		len(response.Additional)))

	// encode response dns message
	resBuffer := Encode(response)

	// send response datagram back to client
	_, err = s.conn.WriteToUDP(resBuffer, clientAddr)
	if err != nil {
		log.Println("failed to send response: ", err)
	}
}
