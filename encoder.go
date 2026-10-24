package godns

import (
	"encoding/binary"
	"io"
	"net"
	"strings"
)

type DnsMessageEncoder struct {
	message DnsMessage
}

func NewDnsMessageEncoder(message DnsMessage) *DnsMessageEncoder {
	return &DnsMessageEncoder{message: message}
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

	_, err := w.Write(b)
	return err
}

func (e *DnsMessageEncoder) writeQuestion(question DnsQuestion, w io.Writer) error {
	// encode qname
	b := []byte{}
	names := strings.Split(question.Qname, ".")
	for _, name := range names {
		length := len(name)
		if length == 0 {
			continue
		}
		label := []byte(name)
		b = append(b, byte(length))
		b = append(b, label...)
	}
	b = append(b, byte(0))
	b = binary.BigEndian.AppendUint16(b, uint16(question.Qtype))
	b = binary.BigEndian.AppendUint16(b, uint16(question.Qclass))
	_, err := w.Write(b)
	return err
}

func (e *DnsMessageEncoder) writeRecord(record DnsRecord, w io.Writer) error {
	b := []byte{}
	b = append(b, e.encodeName(record.Name)...)
	b = binary.BigEndian.AppendUint16(b, uint16(record.Type))
	b = binary.BigEndian.AppendUint16(b, uint16(record.Class))
	b = binary.BigEndian.AppendUint32(b, uint32(record.TTL))
	if record.Type == TypeA {
		encodedRData := encodeIpv4(record.RData)
		b = binary.BigEndian.AppendUint16(b, uint16(len(encodedRData)))
		b = append(b, encodedRData...)
	}
	_, err := w.Write(b)
	return err
}

func (e *DnsMessageEncoder) encodeName(name string) []byte {
	b := []byte{}
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
	return b
}

func encodeIpv4(ipv4 string) []byte {
	ip := net.ParseIP(ipv4).To4()
	return []byte(ip)
}
