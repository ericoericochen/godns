package godns

import "math/rand/v2"

func normalizeName(name string) string {
	return ""
}

func generateTransactionID() uint16 {
	return uint16(rand.UintN(1 << 16))
}
