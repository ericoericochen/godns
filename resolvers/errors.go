package resolvers

import "github.com/ericoericochen/godns"

func formatErrorResponse(message godns.DnsMessage) godns.DnsMessage {
	builder := godns.NewDnsMessageBuilder().From(message).Rcode(godns.RcodeFormatError)
	return builder.Build()
}
