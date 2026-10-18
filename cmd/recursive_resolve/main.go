package main

func main() {
	// rootNameServer := "192.203.230.10"
	// resolver := godns.NewRecursiveResolver(rootNameServer)
	// qname := "google.com"
	// qtype := godns.TypeA
	// qclass := godns.ClassIN

	// build query dns message
	// message := godns.NewDNSMessage()
	// id := uint16(rand.UintN(1 << 16))
	// var id uint16 = 2
	// message.Header.ID = id
	// message.Header.QR = false
	// message.Header.Qdcount = 1
	// message.Questions = append(message.Questions, &godns.DNSQuestion{
	// 	Qname:  qname,
	// 	Qtype:  qtype,
	// 	Qclass: qclass,
	// })

	// res, err := resolver.Query(message)
	// res, err := resolver.Resolve(qname, qtype, qclass)
	// fmt.Println("------")
	// fmt.Println(err)

	// fmt.Println(res.Header)
	// fmt.Println(res.Questions[0])
	// fmt.Println(res.Authorities[0])
	// fmt.Println(res)
}
