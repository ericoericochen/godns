package godns

import (
	"encoding/binary"
	"errors"
	"fmt"
)

/**
+---------------------+
|        Header       |
+---------------------+
|       Question      | the question for the name server
+---------------------+
|        Answer       | RRs answering the question
+---------------------+
|      Authority      | RRs pointing toward an authority
+---------------------+
|      Additional     | RRs holding additional information
+---------------------+
*/

type DNSMessage struct {
	Header      *DNSHeader
	Questions   []*DNSQuestion
	Answers     []*DNSRecord
	Authorities []*DNSRecord
	Additionals []*DNSRecord
}

// Header

/**
							   1  1  1  1  1  1
 0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                      ID                       |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    QDCOUNT                    |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    ANCOUNT                    |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    NSCOUNT                    |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    ARCOUNT                    |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
*/

type DNSHeader struct {
	ID                  uint16       // (16-bit) identifier
	QR                  bool         // (1-bit) 0 - query, 1 - response
	Opcode              Opcode       // (4-bit) kind of query 0 - QUERY, 1 - IQUERY, 2 - STATUS
	AuthoritativeAnswer bool         // (1-bit) is the responding nameserver an authority for the domain?
	Truncation          bool         // (1-bit) is the message truncated?
	RecursionDesired    bool         // (1-bit) ask the nameserver to resolve query recursively
	RecursionAvailable  bool         // (1-bit) does nameserver support recursive queries?
	Z                   uint8        // (3-bit) must be set to 0
	ResponseCode        ResponseCode // (4-bit) response code
	Qdcount             uint16       // (16-bit) number of entries in question section
	Ancount             uint16       // (16-bit) number of resource records in answer section
	Nscount             uint16       // (16-bit) number of nameserver resource records in authority section
	Arcount             uint16       // (16-bit) number of resource records in addition records section
}

func (header *DNSHeader) String() string {
	return fmt.Sprintf(`DNSHeader { ID: %v, QR: %v, OPCODE: %v }`, header.ID, header.QR, header.Opcode)
}

// Question Section

type DNSQuestion struct {
	Qname  string
	Qtype  Type
	Qclass Class
}

func (q *DNSQuestion) String() string {
	return fmt.Sprintf(`DNSQuestion { QNAME: %v, QTYPE: %v, QCLASS: %v }`,
		q.Qname, typeNameMap[q.Qtype], classNameMap[q.Qclass])
}

// DNS Resource Record
type DNSRecord struct {
}

// Parsers
var MalformedDNSMessageError = errors.New("Malformed DNS Message")

func NewDnsMessageFromBuffer(buffer []byte) (*DNSMessage, error) {
	// parse dns header
	// 12 bytes
	id := binary.BigEndian.Uint16(buffer[:2])
	qr := (buffer[2]>>7)&1 == 1
	opcode := (buffer[2] >> 4) & 0b111
	aa := (buffer[2]>>2)&1 == 1
	tc := (buffer[2]>>1)&1 == 1
	rd := (buffer[2] & 1) == 1
	ra := (buffer[3]>>7)&1 == 1
	z := (buffer[3] >> 4) & 0b111
	rcode := buffer[3] & 0b111
	qdcount := binary.BigEndian.Uint16(buffer[4:6])
	ancount := binary.BigEndian.Uint16(buffer[6:8])
	nscount := binary.BigEndian.Uint16(buffer[8:10])
	arcount := binary.BigEndian.Uint16(buffer[10:12])
	dnsHeader := &DNSHeader{}
	dnsHeader.ID = id
	dnsHeader.QR = qr
	dnsHeader.Opcode = Opcode(opcode)
	dnsHeader.AuthoritativeAnswer = aa
	dnsHeader.Truncation = tc
	dnsHeader.RecursionDesired = rd
	dnsHeader.RecursionAvailable = ra
	dnsHeader.Z = z
	dnsHeader.ResponseCode = ResponseCode(rcode)
	dnsHeader.Qdcount = qdcount
	dnsHeader.Ancount = ancount
	dnsHeader.Nscount = nscount
	dnsHeader.Arcount = arcount

	// parse questions section
	questions := []*DNSQuestion{}
	idx := 12
	for i := 0; i < int(qdcount); i++ {
		// parse the qname
		// qname is represented by a sequence of labels
		// with the format [l1],[byte]xl1,...,[ln],[byte]xln,[0]
		// ending in the zero byte (root)
		qname := ""
		for buffer[idx] != byte(0) {
			length := int(buffer[idx])
			start := idx + 1
			end := start + length
			domain := string(buffer[start:end])
			qname += domain + "."

			// update idx to next length octet or zero byte
			idx = end
		}

		qtype := Type(binary.BigEndian.Uint16(buffer[idx+1 : idx+3]))
		qclass := Class(binary.BigEndian.Uint16(buffer[idx+3 : idx+5]))
		question := &DNSQuestion{
			Qname:  qname,
			Qtype:  qtype,
			Qclass: qclass,
		}
		questions = append(questions, question)
		idx += 5
	}

	// build dns message
	message := &DNSMessage{}
	message.Header = dnsHeader
	message.Questions = questions

	return message, nil
}
