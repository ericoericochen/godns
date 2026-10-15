package godns

import (
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

func NewDNSMessage() *DNSMessage {
	return &DNSMessage{
		Header:      &DNSHeader{},
		Questions:   []*DNSQuestion{},
		Answers:     []*DNSRecord{},
		Authorities: []*DNSRecord{},
		Additionals: []*DNSRecord{},
	}
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

//                                1  1  1  1  1  1
//  0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                                               |
// /                                               /
// /                      NAME                     /
// |                                               |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                      TYPE                     |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                     CLASS                     |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                      TTL                      |
// |                                               |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                   RDLENGTH                    |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--|
// /                     RDATA                     /
// /                                               /
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

type DNSRecord struct {
	Name  string
	Type  Type
	Class Class
	TTL   uint32
	RData string
}

func (r *DNSRecord) String() string {
	return fmt.Sprintf(`DNSRecord { NAME: %v, TYPE: %v, CLASS: %v, TTL: %v, RDATA: %v }`,
		r.Name, typeNameMap[r.Type], classNameMap[r.Class], r.TTL, r.RData)
}

func printByte(b byte) {
	fmt.Printf("%08b\n", b)
}
