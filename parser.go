package godns

import (
	"encoding/binary"
	"errors"
	"fmt"
)

var MalformedDNSMessageError = errors.New("Malformed DNS Message")

type DNSMessageParser struct {
	buffer        []byte // dns message bytes
	idx           int    // cursor pointer
	offsetsToName map[int]string
}

func NewDNSMessageParser(buffer []byte) *DNSMessageParser {
	parser := &DNSMessageParser{}
	parser.buffer = buffer
	parser.idx = 0
	parser.offsetsToName = map[int]string{}
	return parser
}

func (p *DNSMessageParser) Parse(dest *DNSMessage) error {
	buffer := p.buffer

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
	p.idx = 12
	for i := 0; i < int(qdcount); i++ {

		qname, err := p.parseName()
		if err != nil {
			return err
		}

		fmt.Println("qname: ", qname)

		qtype := Type(binary.BigEndian.Uint16(buffer[p.idx : p.idx+2]))
		qclass := Class(binary.BigEndian.Uint16(buffer[p.idx+2 : p.idx+4]))
		question := &DNSQuestion{
			Qname:  qname,
			Qtype:  qtype,
			Qclass: qclass,
		}
		questions = append(questions, question)
		p.idx += 4
	}

	// parse answer, authority, and additional sections
	totalRRs := ancount + nscount + arcount

	for i := 0; i < int(totalRRs); i++ {
		name, err := p.parseName()
		if err != nil {
			return err
		}

		fmt.Println("name: ", name)
		break

	}

	// build dns message
	dest.Header = dnsHeader
	dest.Questions = questions

	return nil
}

// parse the name starting from the current idx
func (p *DNSMessageParser) parseName() (string, error) {
	name := ""
	offset := p.idx
	for p.buffer[p.idx] != byte(0) {
		// is label or message compression pointer
		first := p.buffer[p.idx]
		if isLabelByte(first) {
			length := int(first)
			start := p.idx + 1
			end := start + length
			domain := string(p.buffer[start:end])
			name += domain + "."
			p.idx = end
		} else if isMessagePointerByte(first) {
			p.idx += 1
			offset := binary.BigEndian.Uint16([]byte{
				first & byte(0b00111111),
				p.buffer[p.idx],
			})
			domain, exists := p.offsetsToName[int(offset)]
			if !exists {
				return "", MalformedDNSMessageError
			}
			name += domain
			break
		}
	}

	p.offsetsToName[offset] = name

	// advance idx by 1 to segment after name
	p.idx += 1
	return name, nil
}

func isLabelByte(b byte) bool {
	return (b >> 6) == 0
}

func isMessagePointerByte(b byte) bool {
	return b == byte(0b11000000)
}

// parse the qname
// qname is represented by a sequence of labels
// with the format [l1],[byte]xl1,...,[ln],[byte]xln,[0]
// ending in the zero byte (root)
