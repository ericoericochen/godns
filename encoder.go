package godns

import (
	"encoding/binary"
	"io"
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
