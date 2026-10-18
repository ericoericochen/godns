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

func (b *DnsMessageBuilder) Rcode(rcode Rcode) *DnsMessageBuilder {
	b.message.Header.Rcode = rcode
	return b
}

func (b *DnsMessageBuilder) Build() DnsMessage {
	return *b.message
}
