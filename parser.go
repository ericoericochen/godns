package godns

import (
	"encoding/binary"
	"errors"
	"net/netip"
)

var MalformedDNSMessageError = errors.New("Malformed DNS Message")
var MalformedRDataError = errors.New("Malformed RData")

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
	totalRecords := ancount + nscount + arcount
	records := []*DNSRecord{}

	for i := 0; i < int(totalRecords); i++ {
		name, err := p.parseName()
		if err != nil {
			return err
		}
		recordType := Type(binary.BigEndian.Uint16(buffer[p.idx : p.idx+2]))
		class := Class(binary.BigEndian.Uint16(buffer[p.idx+2 : p.idx+4]))
		ttl := binary.BigEndian.Uint32(buffer[p.idx+4 : p.idx+8])
		rdlength := binary.BigEndian.Uint16(buffer[p.idx+8 : p.idx+10])
		rdataBuffer := buffer[p.idx+10 : p.idx+10+int(rdlength)]

		var rdata string
		if recordType == TypeA {
			rdata, err = parseIPv4Address(rdataBuffer)
		} else if recordType == TypeNS {
			// clean this up later
			rd, end, err := p.parseNameFrom(p.idx + 10)
			if err != nil {
				return err
			}
			// name exceeded rdlength
			if end != p.idx+10+int(rdlength) {
				return MalformedDNSMessageError
			}
			rdata = rd

		} else if recordType == TypeAAAA {
			rdata, err = parseIPv6Address(rdataBuffer)
		}

		if err != nil {
			return err
		}

		record := &DNSRecord{
			Name:  name,
			Type:  recordType,
			Class: class,
			TTL:   ttl,
			RData: rdata,
		}
		records = append(records, record)
		p.idx += 10 + int(rdlength)
	}

	// build dns message
	dest.Header = dnsHeader
	dest.Questions = questions
	dest.Answers = records[:ancount]
	dest.Authorities = records[ancount : ancount+nscount]
	dest.Additionals = records[ancount+nscount:]

	return nil
}

// parse the name starting from the current idx
func (p *DNSMessageParser) parseName() (string, error) {
	name, idx, err := p.parseNameFrom(p.idx)
	if err != nil {
		return "", err
	}
	p.idx = idx
	return name, err
}

// parse the name starting from `idx`
func (p *DNSMessageParser) parseNameFrom(idx int) (string, int, error) {
	name := ""
	initialOffset := idx
	for p.buffer[idx] != byte(0) {
		// is label or message compression pointer
		first := p.buffer[idx]
		if isLabelByte(first) {
			length := int(first)
			start := idx + 1
			end := start + length
			domain := string(p.buffer[start:end])
			name += domain + "."
			idx = end
		} else if isMessagePointerByte(first) {
			idx += 1
			offset := binary.BigEndian.Uint16([]byte{
				first & byte(0b00111111),
				p.buffer[idx],
			})
			domain, exists := p.offsetsToName[int(offset)]
			if exists {
				name += domain
			} else {
				// recursively read name from offset
				domain, _, err := p.parseNameFrom(int(offset))
				if err != nil {
					return "", 0, MalformedDNSMessageError
				}
				name += domain
			}

			break
		}
	}

	p.offsetsToName[initialOffset] = name

	// advance idx by 1 to segment after name
	idx += 1
	return name, idx, nil
}

func isLabelByte(b byte) bool {
	return (b >> 6) == 0
}

func isMessagePointerByte(b byte) bool {
	return b == byte(0b11000000)
}

func parseIPv4Address(buffer []byte) (string, error) {
	if len(buffer) != 4 {
		return "", MalformedRDataError
	}
	ipv4 := netip.AddrFrom4([4]byte(buffer))
	return ipv4.String(), nil
}

func parseIPv6Address(buffer []byte) (string, error) {
	if len(buffer) != 16 {
		return "", MalformedRDataError
	}
	ipv6 := netip.AddrFrom16([16]byte(buffer))
	return ipv6.String(), nil
}

// parse the qname
// qname is represented by a sequence of labels
// with the format [l1],[byte]xl1,...,[ln],[byte]xln,[0]
// ending in the zero byte (root)
