package randid

import (
	"math/rand/v2"
	"unsafe"
)

// Size is the length in bytes of the ID
const Size = 16

// StringLen is the length of the string representation of ID
const StringLen = 22

// ID is our 128-bit random value
type ID [2]uint64

func (id ID) Bytes() [Size]byte {
	var out [Size]byte

	out[0] = byte(id[0])
	out[1] = byte(id[0] >> 8)
	out[2] = byte(id[0] >> 16)
	out[3] = byte(id[0] >> 24)
	out[4] = byte(id[0] >> 32)
	out[5] = byte(id[0] >> 40)
	out[6] = byte(id[0] >> 48)
	out[7] = byte(id[0] >> 56)

	out[8] = byte(id[1])
	out[9] = byte(id[1] >> 8)
	out[10] = byte(id[1] >> 16)
	out[11] = byte(id[1] >> 24)
	out[12] = byte(id[1] >> 32)
	out[13] = byte(id[1] >> 40)
	out[14] = byte(id[1] >> 48)
	out[15] = byte(id[1] >> 56)

	return out
}

// String returns base64 encoding of our ID
func (id ID) String() string {
	var buf [StringLen]byte
	encodeUnrolled(&buf, id)
	// Directly convert the array to a string without an
	// intermediary slice
	return unsafe.String((*byte)(unsafe.Pointer(&buf)), StringLen)
}

// New generates a new random ID
func New() ID {
	return ID{
		rand.Uint64(),
		rand.Uint64(),
	}
}

// vendoring in a bit simpler variant of base64 url encoding
// that removes some extra branches and removes the concept of padding
// We're working with a fixed size and fixed character set.
var encodeMap = [64]byte{
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P',
	'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f',
	'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v',
	'w', 'x', 'y', 'z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-', '_',
}

func encodeUnrolled(dst *[StringLen]byte, src ID) {
	// Work directly from the two 64-bit words instead of repeatedly
	// reconstructing 24-bit groups byte-by-byte.
	lo := src[0]
	hi := src[1]

	dst[0] = encodeMap[(lo>>2)&0x3F]
	dst[1] = encodeMap[((lo&0x03)<<4)|((lo>>12)&0x0F)]
	dst[2] = encodeMap[((lo>>6)&0x3C)|((lo>>22)&0x03)]
	dst[3] = encodeMap[(lo>>16)&0x3F]

	dst[4] = encodeMap[(lo>>26)&0x3F]
	dst[5] = encodeMap[((lo>>20)&0x30)|((lo>>36)&0x0F)]
	dst[6] = encodeMap[((lo>>30)&0x3C)|((lo>>46)&0x03)]
	dst[7] = encodeMap[(lo>>40)&0x3F]

	dst[8] = encodeMap[(lo>>50)&0x3F]
	dst[9] = encodeMap[((lo>>48)&0x03)<<4|((lo>>60)&0x0F)]
	dst[10] = encodeMap[((lo>>54)&0x3C)|((hi>>6)&0x03)]
	dst[11] = encodeMap[hi&0x3F]

	dst[12] = encodeMap[(hi>>10)&0x3F]
	dst[13] = encodeMap[((hi>>8)&0x03)<<4|((hi>>20)&0x0F)]
	dst[14] = encodeMap[((hi>>14)&0x3C)|((hi>>30)&0x03)]
	dst[15] = encodeMap[(hi>>24)&0x3F]

	dst[16] = encodeMap[(hi>>34)&0x3F]
	dst[17] = encodeMap[((hi>>32)&0x03)<<4|((hi>>44)&0x0F)]
	dst[18] = encodeMap[((hi>>38)&0x3C)|((hi>>54)&0x03)]
	dst[19] = encodeMap[(hi>>48)&0x3F]

	dst[20] = encodeMap[(hi>>58)&0x3F]
	dst[21] = encodeMap[((hi>>56)&0x03)<<4]
}
