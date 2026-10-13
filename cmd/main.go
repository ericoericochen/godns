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
	data, err := os.ReadFile("response.bin")
	if err != nil {
		panic(err)
	}

	// parse into DNS message
	message := &godns.DNSMessage{}
	parser := godns.NewDNSMessageParser(data)
	err = parser.Parse(message)
	// message, err := godns.NewDnsMessageFromBuffer(data)
	// fmt.Println(message.Header)

	fmt.Println(message.Questions[0])
}
