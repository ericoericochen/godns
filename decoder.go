package godns

import (
	"encoding/binary"
	"errors"
	"net/netip"
)

var MalformedDNSMessageError = errors.New("Malformed DNS Message")
var MalformedRDataError = errors.New("Malformed RData")
var ReadExceedsBufferLengthError = errors.New("Read Exceeds Buffer Length")

type DnsMessageDecoder struct {
	buffer        []byte         // dns message bytes
	idx           int            // cursor pointer
	offsetsToName map[int]string // map of offset idx in buffer to name
}

func NewDnsMessageDecoder(buffer []byte) *DnsMessageDecoder {
	parser := &DnsMessageDecoder{}
	parser.buffer = buffer
	parser.idx = 0
	parser.offsetsToName = map[int]string{}
	return parser
}

func Decode(buf []byte) (DnsMessage, error) {
	var message DnsMessage
	decoder := NewDnsMessageDecoder(buf)
	err := decoder.Decode(&message)
	return message, err
}

// cursor manipulation and read methods

// get the length of the buffer
func (p *DnsMessageDecoder) size() int {
	return len(p.buffer)
}

// read a range of bytes, returns the byte slice.
// returns an error if the range exceeds the size of the buffer.
func (p *DnsMessageDecoder) readRange(start, length int) ([]byte, error) {
	if start+length > p.size() {
		return nil, ReadExceedsBufferLengthError
	}
	return p.buffer[start : start+length], nil
}

// read a byte at a specified index
// returns an error if the idx is beyond the buffer size
func (p *DnsMessageDecoder) readAt(idx int) (byte, error) {
	if idx >= p.size() {
		return 0, ReadExceedsBufferLengthError
	}
	return p.buffer[idx], nil
}

// set the current idx position
func (p *DnsMessageDecoder) seek(idx int) {
	p.idx = idx
}

func (p *DnsMessageDecoder) step(steps int) {
	p.idx += steps
}

func (p *DnsMessageDecoder) readUint16() (uint16, error) {
	bytes, err := p.readRange(p.idx, 2)
	if err != nil {
		return 0, err
	}
	p.step(2)
	return binary.BigEndian.Uint16(bytes), nil
}

func (p *DnsMessageDecoder) readUint32() (uint32, error) {
	bytes, err := p.readRange(p.idx, 4)
	if err != nil {
		return 0, err
	}
	p.step(4)
	return binary.BigEndian.Uint32(bytes), nil
}

func (p *DnsMessageDecoder) Decode(dest *DnsMessage) error {
	// parse dns header 12 bytes

	// parse transaction id
	idBytes, err := p.readRange(0, 2)
	if err != nil {
		return err
	}
	id := binary.BigEndian.Uint16(idBytes)

	// parse qr prefix byte
	qrPrefixByte, err := p.readAt(2)
	if err != nil {
		return err
	}

	qr := (qrPrefixByte>>7)&1 == 1
	opcode := (qrPrefixByte >> 4) & 0b111
	aa := (qrPrefixByte>>2)&1 == 1
	tc := (qrPrefixByte>>1)&1 == 1
	rd := (qrPrefixByte & 1) == 1

	// parse ra prefix byte
	raPrefixByte, err := p.readAt(3)
	if err != nil {
		return err
	}

	ra := (raPrefixByte>>7)&1 == 1
	z := (raPrefixByte >> 4) & 0b111
	rcode := raPrefixByte & 0b111
	p.seek(4)

	// parse qdcount, ancount, nscount, arcount
	qdcount, err := p.readUint16()
	if err != nil {
		return err
	}
	ancount, err := p.readUint16()
	if err != nil {
		return err
	}
	nscount, err := p.readUint16()
	if err != nil {
		return err
	}
	arcount, err := p.readUint16()
	if err != nil {
		return err
	}

	// assemble header
	dnsHeader := DnsHeader{}
	dnsHeader.ID = id
	dnsHeader.QR = qr
	dnsHeader.Opcode = Opcode(opcode)
	dnsHeader.AuthoritativeAnswer = aa
	dnsHeader.Truncation = tc
	dnsHeader.RecursionDesired = rd
	dnsHeader.RecursionAvailable = ra
	dnsHeader.Z = z
	dnsHeader.Rcode = Rcode(rcode)
	dnsHeader.Qdcount = qdcount
	dnsHeader.Ancount = ancount
	dnsHeader.Nscount = nscount
	dnsHeader.Arcount = arcount

	// parse questions section
	questions := []DnsQuestion{}
	for i := 0; i < int(qdcount); i++ {
		qname, err := p.parseName()
		if err != nil {
			return err
		}

		// parse qtype
		qtype16, err := p.readUint16()
		if err != nil {
			return err
		}
		qtype := Type(qtype16)

		// parse qclass
		qclass16, err := p.readUint16()
		if err != nil {
			return err
		}
		qclass := Class(qclass16)

		// assemble question = qname, qtype, qclass
		question := DnsQuestion{
			Qname:  qname,
			Qtype:  qtype,
			Qclass: qclass,
		}
		questions = append(questions, question)
	}

	// parse answer, authority, and additional sections
	totalRecords := ancount + nscount + arcount
	records := []DnsRecord{}

	for i := 0; i < int(totalRecords); i++ {
		name, err := p.parseName()
		if err != nil {
			return err
		}

		// parse record type
		recordType16, err := p.readUint16()
		if err != nil {
			return err
		}
		recordType := Type(recordType16)

		// parse class
		class16, err := p.readUint16()
		if err != nil {
			return err
		}
		class := Class(class16)

		// parse ttl
		ttl, err := p.readUint32()
		if err != nil {
			return err
		}

		// parse rdlength
		rdlength, err := p.readUint16()
		if err != nil {
			return err
		}

		// parse rdata
		rdataBuffer, err := p.readRange(p.idx, int(rdlength))

		var rdata string
		if recordType == TypeA {
			rdata, err = parseIPv4Address(rdataBuffer)
		} else if recordType == TypeNS {
			rdata, err = p.parseNameRData(int(rdlength))
		} else if recordType == TypeCNAME {
			rdata, err = p.parseNameRData(int(rdlength))
		} else if recordType == TypeAAAA {
			rdata, err = parseIPv6Address(rdataBuffer)
		}

		if err != nil {
			return err
		}

		record := DnsRecord{
			Name:  name,
			Type:  recordType,
			Class: class,
			TTL:   ttl,
			RData: rdata,
		}
		records = append(records, record)
		p.step(int(rdlength))
	}

	// build dns message
	dest.Header = dnsHeader
	dest.Questions = questions
	dest.Answers = records[:ancount]
	dest.Authority = records[ancount : ancount+nscount]
	dest.Additional = records[ancount+nscount:]

	return nil
}

// parse the name starting from the current idx
func (p *DnsMessageDecoder) parseName() (string, error) {
	name, idx, err := p.parseNameFrom(p.idx)
	if err != nil {
		return "", err
	}
	p.seek(idx)
	return name, err
}

// parse the name starting from `idx`
func (p *DnsMessageDecoder) parseNameFrom(idx int) (string, int, error) {
	name := ""
	initialOffset := idx // save the initial idx so we can cache the offset -> name mapping
	end := -1            // end is unset to -1
	jumps := 0

	for {
		// prevent infinite message ptr loops
		if jumps == 5 {
			return "", 0, MalformedRDataError
		}

		// read the byte, break if 0 byte (end of name)
		b, err := p.readAt(idx)
		if err != nil {
			return "", 0, err
		}
		if b == byte(0) {
			if end == -1 {
				// advance idx by 1 to get the idx of the byte
				// after the 0 byte
				end = idx + 1
			}
			break
		}

		// is label or message compression pointer
		if isLabelByte(b) {
			// parse <length><char x length> and add to name
			length := int(b)
			start := idx + 1
			domainBytes, err := p.readRange(start, length)
			if err != nil {
				return "", 0, err
			}
			domain := string(domainBytes)
			name += domain + "."
			// move idx to next label to parse
			idx = start + length
		} else if isMessagePointerByte(b) {
			// read the first 2 bytes
			idx += 1
			nb, err := p.readAt(idx)
			if err != nil {
				return "", 0, err
			}
			// message ptr marks the end of name bytes
			if end == -1 {
				// move idx by 1 pass the last message ptr byte
				end = idx + 1
			}
			// set the first 11 of first byte to 00 and parse offset
			offset := binary.BigEndian.Uint16([]byte{
				b & byte(0b00111111), nb,
			})
			// check if we have cached the name at the offset
			// if not we have to jump to the offset and get the name
			domain, exists := p.offsetsToName[int(offset)]
			if exists {
				name += domain
				break
			} else {
				idx = int(offset)
				jumps += 1
				continue
			}
		}
	}

	// cache the offset -> name finding
	p.offsetsToName[initialOffset] = name
	return name, end, nil
}

func (p *DnsMessageDecoder) parseNameRData(rdlength int) (string, error) {
	start := p.idx
	rd, end, err := p.parseNameFrom(start)
	if err != nil {
		return "", err
	}
	if end != start+rdlength {
		return "", MalformedDNSMessageError
	}
	return rd, nil
}

func isLabelByte(b byte) bool {
	return (b >> 6) == 0
}

// message pointer bytes start with 11
func isMessagePointerByte(b byte) bool {
	return b&byte(0b11000000) == 0b11000000
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
