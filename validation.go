package godns

func ValidateQueryDNSMessage(message *DnsMessage) bool {
	if message.Header.QR {
		return false
	}
	if int(message.Header.Qdcount) != len(message.Questions) {
		return false
	}
	return true
}

// func IsValidDNSMessage(message *DNSMessage) bool {
// 	// queries can not have resource records
// 	if !message.Header.QR {
// 		if message.Header.Ancount > 0 || message.Header.Nscount > 0 || message.Header.Arcount > 0 {
// 			return false
// 		}

// 		if len(message.Answers) > 0 || len(message.Authorities) > 0 || len(message.Additionals) > 0 {
// 			return false
// 		}
// 	}

// 	// header matches number of questions and resource records
// 	if int(message.Header.Qdcount) != len(message.Questions) {
// 		return false
// 	}
// 	if int(message.Header.Ancount) != len(message.Answers) {
// 		return false
// 	}
// 	if int(message.Header.Nscount) != len(message.Authorities) {
// 		return false
// 	}
// 	if int(message.Header.Arcount) != len(message.Additionals) {
// 		return false
// 	}

// 	return true
// }
