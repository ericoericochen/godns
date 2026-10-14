package godns

type StubResolver struct {
	server string
}

func NewStubResolver(server string) *StubResolver {
	return &StubResolver{server: server}
}

func (r *StubResolver) Query(message *DNSMessage) (*DNSMessage, error) {
	return nil, nil
}
