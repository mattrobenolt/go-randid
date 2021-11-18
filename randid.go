package randid

import "crypto/rand"

// Size is the length in bytes of the ID
const Size = 16

// StringLen is the length of the string representation of ID
const StringLen = 22

// ID is our 16 byte random value
type ID [Size]byte

// hook for tests to stub in a predictable random
var randReader = rand.Read

// String returns base64 encoding of our ID
func (id ID) String() string {
	var buf [StringLen]byte
	encode(buf[:], id)
	return string(buf[:])
}

// New generates a new random ID
func New() (id ID) {
	if _, err := randReader(id[:]); err != nil {
		panic(err)
	}
	return
}

// vendoringing in a bit simpler variant of base64 url encoding
// that removes some extra branches and removes the concept of padding
// We're working with a fixed size and fixed character set.
var encodeMap = [64]byte{
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P',
	'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f',
	'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v',
	'w', 'x', 'y', 'z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-', '_',
}

func encode(dst []byte, src ID) {
	di, si, val := 0, 0, uint(0)
	n := Size - 1
	for si < n {
		// Convert 3x 8bit source bytes into 4 bytes
		val = uint(src[si+0])<<16 | uint(src[si+1])<<8 | uint(src[si+2])

		dst[di+0] = encodeMap[val>>18&0x3F]
		dst[di+1] = encodeMap[val>>12&0x3F]
		dst[di+2] = encodeMap[val>>6&0x3F]
		dst[di+3] = encodeMap[val&0x3F]

		si += 3
		di += 4
	}

	val = uint(src[si+0]) << 16
	dst[di+0] = encodeMap[val>>18&0x3F]
	dst[di+1] = encodeMap[val>>12&0x3F]
}
