package godns_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/ericoericochen/godns"
)

func TestDecodeResponseFixture(t *testing.T) {
	// Expected messages from the corresponding dig output in fixtures/*.txt.
	fixtures := []struct {
		name string
		want godns.DnsMessage
	}{
		{
			name: "www.google.com",
			want: godns.DnsMessage{
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
			},
		},
		{
			name: "en.wikipedia.org",
			want: godns.DnsMessage{
				Header: godns.DnsHeader{
					ID:                  35270,
					QR:                  true,
					Opcode:              godns.OpcodeQUERY,
					AuthoritativeAnswer: false,
					Truncation:          false,
					RecursionDesired:    true,
					RecursionAvailable:  true,
					Z:                   0,
					Rcode:               godns.RcodeNoError,
					Qdcount:             1,
					Ancount:             2,
					Nscount:             0,
					Arcount:             0,
				},
				Questions: []godns.DnsQuestion{
					{Qname: "en.wikipedia.org.", Qtype: godns.TypeA, Qclass: godns.ClassIN},
				},
				Answers: []godns.DnsRecord{
					{Name: "en.wikipedia.org.", Type: godns.TypeCNAME, Class: godns.ClassIN, TTL: 86181, RData: "dyna.wikimedia.org."},
					{Name: "dyna.wikimedia.org.", Type: godns.TypeA, Class: godns.ClassIN, TTL: 141, RData: "208.80.153.224"},
				},
				Authority:  []godns.DnsRecord{},
				Additional: []godns.DnsRecord{},
			},
		},
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

	for _, fixture := range fixtures {
		t.Run(fixture.name, func(t *testing.T) {
			data, err := os.ReadFile("fixtures/" + fixture.name + ".bin")
			if err != nil {
				t.Fatalf("read response fixture: %v", err)
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					got, err := tt.decode(data)
					if err != nil {
						t.Fatalf("decode response fixture: %v", err)
					}
					if !reflect.DeepEqual(got, fixture.want) {
						t.Errorf("decoded message does not match fixtures/%s.txt:\ngot:  %#v\nwant: %#v", fixture.name, got, fixture.want)
					}
				})
			}
		})
	}
}
