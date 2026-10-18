package godns

import "errors"

type DnsResolver interface {
	WithCache(*DnsCache) DnsResolver
	Resolve(DnsMessage) DnsMessage
}

// validation

var InvalidDnsQueryError = errors.New("InvalidDnsQueryError")

func validateQueryDnsMessage(message DnsMessage) error {
	if message.Header.QR {
		return InvalidDnsQueryError
	}

	if message.Header.Qdcount != 1 {
		return InvalidDnsQueryError
	}

	if len(message.Questions) != 1 {
		return InvalidDnsQueryError
	}

	if message.Header.Ancount != 0 || message.Header.Nscount != 0 || message.Header.Arcount != 0 {
		return InvalidDnsQueryError
	}

	if len(message.Answers) != 0 || len(message.Authority) != 0 || len(message.Additional) != 0 {
		return InvalidDnsQueryError
	}

	return nil
}

// errors
func toResponse(message DnsMessage, rcode Rcode) DnsMessage {
	return NewDnsMessageBuilder().From(message).Response().Rcode(rcode).Build()
}
