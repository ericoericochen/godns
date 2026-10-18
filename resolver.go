package godns

type DnsResolver interface {
	Resolve(*DnsMessage) (*DnsMessage, error)
}
