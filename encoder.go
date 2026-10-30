package godns

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"strings"
)

type DnsMessageEncoder struct {
	message DnsMessage
	cursor  int
	offsets map[string]int
}

func NewDnsMessageEncoder(message DnsMessage) *DnsMessageEncoder {
	return &DnsMessageEncoder{
		message: message,
		cursor:  0,
		offsets: map[string]int{},
	}
}

func Encode(message DnsMessage) []byte {
	var buf bytes.Buffer
	encoder := NewDnsMessageEncoder(message)
	encoder.Write(&buf)
	return buf.Bytes()
}

func (e *DnsMessageEncoder) Write(w io.Writer) error {
	err := e.writeHeader(w)
	if err != nil {
		return err
	}
	for i := 0; i < int(e.message.Header.Qdcount); i++ {
		question := e.message.Questions[i]
		err := e.writeQuestion(question, w)
		if err != nil {
			return err
		}
	}
	records := []DnsRecord{}
	records = append(records, e.message.Answers...)
	records = append(records, e.message.Authority...)
	records = append(records, e.message.Additional...)
	for _, record := range records {
		err := e.writeRecord(record, w)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *DnsMessageEncoder) writeHeader(w io.Writer) error {
	b := []byte{}
	header := e.message.Header
	b = binary.BigEndian.AppendUint16(b, header.ID)

	// qr, opcode, aa, tc, rd
	qrByte := byte(0)
	if header.QR {
		qrByte |= 1 << 7
	}

	// encode opcode
	opcode := byte(header.Opcode)
	qrByte |= opcode << 3

	// encode ra
	if header.RecursionDesired {
		qrByte |= 1
	}

	b = append(b, qrByte)

	// ra, z, rcode
	raByte := byte(0)
	b = append(b, raByte)

	// encode qdcount, ancount, nscount, arcount
	b = binary.BigEndian.AppendUint16(b, header.Qdcount)
	b = binary.BigEndian.AppendUint16(b, header.Ancount)
	b = binary.BigEndian.AppendUint16(b, header.Nscount)
	b = binary.BigEndian.AppendUint16(b, header.Arcount)

	// header takes up 12 bytes, move cursor
	e.cursor += 12

	_, err := w.Write(b)
	return err
}

func (e *DnsMessageEncoder) writeQuestion(question DnsQuestion, w io.Writer) error {
	// encode qname
	b := []byte{}
	b = append(b, e.encodeName(question.Qname)...)
	b = binary.BigEndian.AppendUint16(b, uint16(question.Qtype))
	b = binary.BigEndian.AppendUint16(b, uint16(question.Qclass))
	e.cursor += len(b)
	_, err := w.Write(b)
	return err
}

func (e *DnsMessageEncoder) writeRecord(record DnsRecord, w io.Writer) error {
	b := []byte{}
	b = append(b, e.encodeName(record.Name)...)
	b = binary.BigEndian.AppendUint16(b, uint16(record.Type))
	b = binary.BigEndian.AppendUint16(b, uint16(record.Class))
	b = binary.BigEndian.AppendUint32(b, uint32(record.TTL))

	var encodedRData []byte

	if record.Type == TypeA {
		encodedRData = encodeIpv4(record.RData)
	} else if record.Type == TypeCNAME {
		encodedRData = e.encodeName(record.RData)
	}

	// add the rdlength then the rdata bytes
	b = binary.BigEndian.AppendUint16(b, uint16(len(encodedRData)))
	b = append(b, encodedRData...)

	e.cursor += len(b)
	_, err := w.Write(b)
	return err
}

func (e *DnsMessageEncoder) encodeName(name string) []byte {
	b := []byte{}
	offset, exists := e.offsets[name]
	if exists {
		ptr := 0b11000000_00000000 | uint16(offset)
		b = binary.BigEndian.AppendUint16(b, ptr)
		return b
	}

	names := strings.Split(name, ".")
	for _, part := range names {
		length := len(part)
		if length == 0 {
			continue
		}
		label := []byte(part)
		b = append(b, byte(length))
		b = append(b, label...)
	}
	b = append(b, byte(0))
	e.offsets[name] = e.cursor

	return b
}

func encodeIpv4(ipv4 string) []byte {
	ip := net.ParseIP(ipv4).To4()
	return []byte(ip)
}
