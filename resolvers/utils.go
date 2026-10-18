package resolvers

import (
	"github.com/ericoericochen/godns"
	"github.com/ericoericochen/godns/cache"
)

func cachedValueToResponseDnsMessage(query *godns.DnsMessage, value cache.DnsCacheValue) *godns.DnsMessage {
	// copy over dns query
	message := godns.NewDnsMessage()
	message.Header.ID = query.Header.ID
	message.Header.QR = true
	message.Header.RecursionDesired = query.Header.RecursionDesired
	message.Header.Qdcount = 1
	message.Questions = append(message.Questions, query.Questions[0])

	// set cached value
	message.Header.Ancount = uint16(len(value.Answers))
	message.Header.Nscount = uint16(len(value.Authority))
	message.Header.Arcount = uint16(len(value.Additional))

	message.Answers = append(message.Answers, value.Answers...)
	message.Authority = append(message.Authority, value.Authority...)
	message.Additional = append(message.Additional, value.Additional...)
	return message
}
