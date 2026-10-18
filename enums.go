package godns

// OPCODE
type Opcode uint8

const (
	OpcodeQUERY Opcode = iota
	OpcodeIQUERY
	OpcodeSTATUS
)

// RCODE
type ResponseCode uint8

const (
	ResponseCodeNoError ResponseCode = iota
	ResponseCodeFormatError
	ResponseCodeServerFailure
	ResponseCodeNameError
	ResponseCodeNotImplemented
	ResponseCodeRefused
)

var responseCodeMap map[ResponseCode]string = map[ResponseCode]string{
	ResponseCodeNoError:        "NOERROR",
	ResponseCodeFormatError:    "FORMERR",
	ResponseCodeServerFailure:  "SERVFAIL",
	ResponseCodeNameError:      "NXDOMAIN",
	ResponseCodeNotImplemented: "NOTIMP",
	ResponseCodeRefused:        "REFUSED",
}

// TYPE, QTYPE
type Type uint8

const (
	TypeA    Type = 1
	TypeNS   Type = 2
	TypeAAAA Type = 28
)

var typeNameMap map[Type]string = map[Type]string{
	TypeA:    "A",
	TypeNS:   "NS",
	TypeAAAA: "AAAA",
}

// CLASS, QCLASS
type Class uint8

const (
	ClassIN  Class = 1
	ClassCS        = 2
	ClassCH        = 3
	ClassHS        = 4
	ClassANY       = 255
)

var classNameMap map[Class]string = map[Class]string{
	ClassIN:  "IN",
	ClassCS:  "CS",
	ClassCH:  "CH",
	ClassHS:  "HS",
	ClassANY: "*",
}
