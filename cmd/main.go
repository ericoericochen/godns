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
	message := &godns.DnsMessage{}
	parser := godns.NewDNSMessageParser(data)
	err = parser.Parse(message)

	for _, record := range message.Answers {
		fmt.Println(record)
	}
	fmt.Println("---")
	for _, record := range message.Additionals {
		fmt.Println(record)
	}
}
