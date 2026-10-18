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

type DnsMessage struct {
	Header     *DnsHeader
	Questions  []*DnsQuestion
	Answers    []*DnsRecord
	Authority  []*DnsRecord
	Additional []*DnsRecord
}

func NewDnsMessage() *DnsMessage {
	header := &DnsHeader{}
	header.ID = generateTransactionID()
	return &DnsMessage{
		Header:     header,
		Questions:  []*DnsQuestion{},
		Answers:    []*DnsRecord{},
		Authority:  []*DnsRecord{},
		Additional: []*DnsRecord{},
	}
}

func NewDnsQuery(qname string, qtype Type, qclass Class, rd bool) *DnsMessage {
	message := NewDnsMessage()
	message.Header.Qdcount = 1
	message.Header.RecursionDesired = rd
	question := &DnsQuestion{
		Qname:  qname,
		Qtype:  qtype,
		Qclass: qclass,
	}
	message.Questions = append(message.Questions, question)
	return message
}

func (m *DnsMessage) String() string {
	s := "DnsMessage {\n"
	// header
	s += "\t" + m.Header.String() + "\n"
	// questions
	s += "\t" + "Questions {" + "\n"
	for _, question := range m.Questions {
		s += "\t\t" + question.String() + "\n"
	}
	s += "\t}\n"
	// answers
	s += "\tAnswers{\n"
	for _, answer := range m.Answers {
		s += "\t\t" + answer.String() + "\n"
	}
	s += "\t}\n"
	// authorities
	s += "\tAuthorities {\n"
	for _, record := range m.Authority {
		s += "\t\t" + record.String() + "\n"
	}
	s += "\t}\n"
	// additionals
	s += "\tAdditionals {\n"
	for _, record := range m.Additional {
		s += "\t\t" + record.String() + "\n"
	}
	s += "\t}\n"

	s += "}"
	return s
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

type DnsHeader struct {
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

func (header *DnsHeader) String() string {
	return fmt.Sprintf(`Header { ID: %v, QR: %v, OPCODE: %v, AA: %v, TC: %v, RD: %v, RA: %v, Z: %v, RCODE: %v, QDCOUNT: %v, ANCOUNT: %v, NSCOUNT: %v, ARCOUNT: %v }`,
		header.ID, header.QR, header.Opcode, header.AuthoritativeAnswer, header.Truncation,
		header.RecursionDesired, header.RecursionAvailable, header.Z, responseCodeMap[header.ResponseCode],
		header.Qdcount, header.Ancount, header.Nscount, header.Arcount)
}

// Question Section

//                                1  1  1  1  1  1
//  0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                                               |
// /                     QNAME                     /
// /                                               /
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                     QTYPE                     |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
// |                     QCLASS                    |
// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

type DnsQuestion struct {
	Qname  string
	Qtype  Type
	Qclass Class
}

func (q *DnsQuestion) String() string {
	return fmt.Sprintf(`Question { QNAME: %v, QTYPE: %v, QCLASS: %v }`,
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

type DnsRecord struct {
	Name  string
	Type  Type
	Class Class
	TTL   uint32
	RData string
}

func (r *DnsRecord) String() string {
	return fmt.Sprintf(`Record { NAME: %v, TYPE: %v, CLASS: %v, TTL: %v, RDATA: %v }`,
		r.Name, typeNameMap[r.Type], classNameMap[r.Class], r.TTL, r.RData)
}

func printByte(b byte) {
	fmt.Printf("%08b\n", b)
}
