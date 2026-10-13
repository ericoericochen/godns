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
	ResponseCodeNameError
	ResponseCodeNotImplemented
	ResponseCodeRefused
)

// TYPE, QTYPE
type Type uint8

const (
	TypeA  Type = 1
	TypeNS      = 2
)

var typeNameMap map[Type]string = map[Type]string{
	TypeA:  "A",
	TypeNS: "NS",
}

// CLASS, QCLASS
type Class uint8

const (
	ClassIN  = 1
	ClassCS  = 2
	ClassCH  = 3
	ClassHS  = 4
	ClassANY = 255
)

var classNameMap map[Class]string = map[Class]string{
	ClassIN:  "IN",
	ClassCS:  "CS",
	ClassCH:  "CH",
	ClassHS:  "HS",
	ClassANY: "*",
}
