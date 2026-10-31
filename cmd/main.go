package main

import (
	"fmt"
	"os"

	"github.com/ericoericochen/godns"
)

func printByte(b byte) {
	fmt.Printf("%08b\n", b)
}

func main() {
	// data, err := os.ReadFile("fixtures/response.bin")
	// data, err := os.ReadFile("fixtures/www.google.com.bin")
	data, err := os.ReadFile("fixtures/en.wikipedia.org.bin")
	if err != nil {
		panic(err)
	}

	// parse into DNS message
	message := &godns.DnsMessage{}
	parser := godns.NewDnsMessageDecoder(data)
	err = parser.Decode(message)

	fmt.Println(message)
}
