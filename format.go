package godns

import (
	"fmt"
	"strings"
)

func DigFormat(message *DnsMessage) string {
	var b strings.Builder
	header := message.Header

	opcodes := map[Opcode]string{
		OpcodeQUERY:  "QUERY",
		OpcodeIQUERY: "IQUERY",
		OpcodeSTATUS: "STATUS",
	}

	flags := []string{}
	if header.QR {
		flags = append(flags, "qr")
	}
	if header.AuthoritativeAnswer {
		flags = append(flags, "aa")
	}
	if header.Truncation {
		flags = append(flags, "tc")
	}
	if header.RecursionDesired {
		flags = append(flags, "rd")
	}
	if header.RecursionAvailable {
		flags = append(flags, "ra")
	}

	fmt.Fprintf(&b, ";; ->>HEADER<<- opcode: %s, status: %s, id: %d\n",
		opcodes[header.Opcode], responseCodeMap[header.ResponseCode], header.ID)
	fmt.Fprintf(&b, ";; flags: %s; QUERY: %d, ANSWER: %d, AUTHORITY: %d, ADDITIONAL: %d\n",
		strings.Join(flags, " "), header.Qdcount, header.Ancount, header.Nscount, header.Arcount)

	if len(message.Questions) > 0 {
		b.WriteString("\n;; QUESTION SECTION:\n")
		for _, question := range message.Questions {
			fmt.Fprintf(&b, ";%s\t%s\t%s\n",
				question.Qname, classNameMap[question.Qclass], typeNameMap[question.Qtype])
		}
	}

	sections := []struct {
		name    string
		records []*DnsRecord
	}{
		{"ANSWER", message.Answers},
		{"AUTHORITY", message.Authority},
		{"ADDITIONAL", message.Additional},
	}

	for _, section := range sections {
		if len(section.records) == 0 {
			continue
		}
		fmt.Fprintf(&b, "\n;; %s SECTION:\n", section.name)
		for _, record := range section.records {
			fmt.Fprintf(&b, "%s\t%d\t%s\t%s\t%s\n",
				record.Name, record.TTL, classNameMap[record.Class], typeNameMap[record.Type], record.RData)
		}
	}

	return strings.TrimSuffix(b.String(), "\n")
}
