package godns

// OPCODE
type Opcode uint8

const (
	OpcodeQUERY Opcode = iota
	OpcodeIQUERY
	OpcodeSTATUS
)

// RCODE
type Rcode uint8

const (
	RcodeNoError Rcode = iota
	RcodeFormatError
	RcodeServerFailure
	RcodeNameError
	RcodeNotImplemented
	RcodeRefused
)

var rcodeMap map[Rcode]string = map[Rcode]string{
	RcodeNoError:        "NOERROR",
	RcodeFormatError:    "FORMERR",
	RcodeServerFailure:  "SERVFAIL",
	RcodeNameError:      "NXDOMAIN",
	RcodeNotImplemented: "NOTIMP",
	RcodeRefused:        "REFUSED",
}

// TYPE, QTYPE
type Type uint8

const (
	TypeA     Type = 1
	TypeNS    Type = 2
	TypeCNAME Type = 5
	TypeAAAA  Type = 28
)

var typeNameMap map[Type]string = map[Type]string{
	TypeA:     "A",
	TypeNS:    "NS",
	TypeCNAME: "CNAME",
	TypeAAAA:  "AAAA",
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
