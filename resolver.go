package godns

// Dns Resolver receives a query and returns a response.
// Both the query and response are Dns Messages.
type DnsResolver interface {
	Resolve(DnsMessage) DnsMessage
}
