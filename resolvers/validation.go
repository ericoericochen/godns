package resolvers

import (
	"errors"

	"github.com/ericoericochen/godns"
)

var InvalidDnsQueryError = errors.New("InvalidDnsQueryError")

func validateQueryDnsMessage(message godns.DnsMessage) error {
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
