package godns

type DnsMessageBuilder struct {
	message *DnsMessage
}

func NewDnsMessageBuilder() *DnsMessageBuilder {
	message := NewDnsMessage()
	return &DnsMessageBuilder{message: &message}
}

func (b *DnsMessageBuilder) From(message DnsMessage) *DnsMessageBuilder {
	b.message = &message
	return b
}

func (b *DnsMessageBuilder) Query(qname string, qtype Type, qclass Class) *DnsMessageBuilder {
	message := b.message
	question := DnsQuestion{Qname: qname, Qtype: qtype, Qclass: qclass}
	message.Header.QR = false
	message.Header.Qdcount = 1
	message.Questions = []DnsQuestion{question}
	message.Answers = []DnsRecord{}
	message.Authority = []DnsRecord{}
	message.Additional = []DnsRecord{}
	return b
}

func (b *DnsMessageBuilder) RecursionDesired(rd bool) *DnsMessageBuilder {
	b.message.Header.RecursionDesired = rd
	return b
}

func (b *DnsMessageBuilder) RecursionAvailable(ra bool) *DnsMessageBuilder {
	b.message.Header.RecursionAvailable = ra
	return b
}

func (b *DnsMessageBuilder) Response() *DnsMessageBuilder {
	b.message.Header.QR = true
	return b
}

func (b *DnsMessageBuilder) Rcode(rcode Rcode) *DnsMessageBuilder {
	b.message.Header.Rcode = rcode
	return b
}

func (b *DnsMessageBuilder) Answers(answers []DnsRecord) *DnsMessageBuilder {
	b.message.Answers = append(b.message.Answers, answers...)
	b.message.Header.Ancount = uint16(len(b.message.Answers))
	return b
}

func (b *DnsMessageBuilder) Authority(authority []DnsRecord) *DnsMessageBuilder {
	b.message.Authority = append(b.message.Authority, authority...)
	b.message.Header.Nscount = uint16(len(b.message.Authority))
	return b
}

func (b *DnsMessageBuilder) Additional(additional []DnsRecord) *DnsMessageBuilder {
	b.message.Additional = append(b.message.Additional, additional...)
	b.message.Header.Arcount = uint16(len(b.message.Additional))
	return b
}

func (b *DnsMessageBuilder) Build() DnsMessage {
	return *b.message
}
