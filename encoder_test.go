package godns_test

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"github.com/ericoericochen/godns"
)

func TestEncodeResponseFixtureRoundTrip(t *testing.T) {
	tests := []struct {
		name   string
		encode func(godns.DnsMessage) ([]byte, error)
	}{
		{
			name: "Encode",
			encode: func(message godns.DnsMessage) ([]byte, error) {
				return godns.Encode(message), nil
			},
		},
		{
			name: "DnsMessageEncoder.Write",
			encode: func(message godns.DnsMessage) ([]byte, error) {
				var buf bytes.Buffer
				err := godns.NewDnsMessageEncoder(message).Write(&buf)
				return buf.Bytes(), err
			},
		},
	}

	for _, name := range []string{"www.google.com", "en.wikipedia.org"} {
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile("fixtures/" + name + ".bin")
			if err != nil {
				t.Fatalf("read response fixture: %v", err)
			}
			want, err := godns.Decode(data)
			if err != nil {
				t.Fatalf("decode response fixture: %v", err)
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					encoded, err := tt.encode(want)
					if err != nil {
						t.Fatalf("encode decoded message: %v", err)
					}
					got, err := godns.Decode(encoded)
					if err != nil {
						t.Fatalf("decode encoded message: %v", err)
					}
					if !reflect.DeepEqual(got, want) {
						t.Errorf("round-trip message differs from original:\ngot:  %#v\nwant: %#v", got, want)
					}
				})
			}
		})
	}
}
