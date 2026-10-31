package main

import (
	"fmt"

	"github.com/ericoericochen/godns"
)

func printByte(b byte) {
	fmt.Printf("%08b\n", b)
}

func main() {
	var err error
	// data, err := os.ReadFile("fixtures/response.bin")
	// data, err := os.ReadFile("fixtures/www.google.com.bin")
	// data, err := os.ReadFile("fixtures/en.wikipedia.org.bin")
	// if err != nil {
	// 	panic(err)
	// }

	data := []byte{
		110, 60, 1, 32, 0, 1, 0, 0, 0, 0, 0, 0,
		3, 119, 119, 119, 6, 103, 111, 111, 103, 108, 101,
		3, 99, 111, 109, 0, 0, 1, 0, 1,
	}

	// parse into DNS message
	message := &godns.DnsMessage{}
	parser := godns.NewDnsMessageDecoder(data)
	err = parser.Decode(message)

	fmt.Println(message)

	fmt.Println("err: ", err)
}
