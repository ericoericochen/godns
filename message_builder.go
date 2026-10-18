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

func (b *DnsMessageBuilder) Rcode(rcode Rcode) *DnsMessageBuilder {
	b.message.Header.Rcode = rcode
	return b
}

func (b *DnsMessageBuilder) Build() DnsMessage {
	return *b.message
}
