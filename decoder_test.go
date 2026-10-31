package godns_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/ericoericochen/godns"
)

func TestDecodeResponseFixture(t *testing.T) {
	data, err := os.ReadFile("fixtures/response.bin")
	if err != nil {
		t.Fatalf("read response fixture: %v", err)
	}

	// Expected message from the dig output in fixtures/response.txt.
	want := godns.DnsMessage{
		Header: godns.DnsHeader{
			ID:                  945,
			QR:                  true,
			Opcode:              godns.OpcodeQUERY,
			AuthoritativeAnswer: false,
			Truncation:          false,
			RecursionDesired:    true,
			RecursionAvailable:  true,
			Z:                   0,
			Rcode:               godns.RcodeNoError,
			Qdcount:             1,
			Ancount:             8,
			Nscount:             0,
			Arcount:             0,
		},
		Questions: []godns.DnsQuestion{
			{Qname: "www.google.com.", Qtype: godns.TypeA, Qclass: godns.ClassIN},
		},
		Answers: []godns.DnsRecord{
			{Name: "www.google.com.", Type: godns.TypeA, Class: godns.ClassIN, TTL: 285, RData: "142.251.151.119"},
			{Name: "www.google.com.", Type: godns.TypeA, Class: godns.ClassIN, TTL: 285, RData: "142.251.155.119"},
			{Name: "www.google.com.", Type: godns.TypeA, Class: godns.ClassIN, TTL: 285, RData: "142.251.150.119"},
			{Name: "www.google.com.", Type: godns.TypeA, Class: godns.ClassIN, TTL: 285, RData: "142.251.152.119"},
			{Name: "www.google.com.", Type: godns.TypeA, Class: godns.ClassIN, TTL: 285, RData: "142.251.156.119"},
			{Name: "www.google.com.", Type: godns.TypeA, Class: godns.ClassIN, TTL: 285, RData: "142.251.154.119"},
			{Name: "www.google.com.", Type: godns.TypeA, Class: godns.ClassIN, TTL: 285, RData: "142.251.157.119"},
			{Name: "www.google.com.", Type: godns.TypeA, Class: godns.ClassIN, TTL: 285, RData: "142.251.153.119"},
		},
		Authority:  []godns.DnsRecord{},
		Additional: []godns.DnsRecord{},
	}

	tests := []struct {
		name   string
		decode func([]byte) (godns.DnsMessage, error)
	}{
		{name: "Decode", decode: godns.Decode},
		{
			name: "DnsMessageDecoder.Decode",
			decode: func(data []byte) (godns.DnsMessage, error) {
				var message godns.DnsMessage
				err := godns.NewDnsMessageDecoder(data).Decode(&message)
				return message, err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.decode(data)
			if err != nil {
				t.Fatalf("decode response fixture: %v", err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("decoded message does not match fixtures/response.txt:\ngot:  %#v\nwant: %#v", got, want)
			}
		})
	}
}
