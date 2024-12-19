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
	// XXX: this is manually unrolled and performance is
	// a few percentage faster, likely the only optimization
	// that can be accomplished here is through ASM and utilizing
	// vectors and probably encode this entire thing in one
	// vector CPU instruction.
	var val uint32

	// Convert 3x 8bit source bytes into 4 bytes
	val = (uint32(byte(src[(0+0)/8]>>((0+0)%8*8)))<<16 |
		uint32(byte(src[(0+1)/8]>>((0+1)%8*8)))<<8 |
		uint32(byte(src[(0+2)/8]>>((0+2)%8*8)))<<0)

	dst[0+0] = encodeMap[val>>18&0x3F]
	dst[0+1] = encodeMap[val>>12&0x3F]
	dst[0+2] = encodeMap[val>>6&0x3F]
	dst[0+3] = encodeMap[val&0x3F]

	// Convert 3x 8bit source bytes into 4 bytes
	val = (uint32(byte(src[(3+0)/8]>>((3+0)%8*8)))<<16 |
		uint32(byte(src[(3+1)/8]>>((3+1)%8*8)))<<8 |
		uint32(byte(src[(3+2)/8]>>((3+2)%8*8)))<<0)

	dst[4+0] = encodeMap[val>>18&0x3F]
	dst[4+1] = encodeMap[val>>12&0x3F]
	dst[4+2] = encodeMap[val>>6&0x3F]
	dst[4+3] = encodeMap[val&0x3F]

	// Convert 3x 8bit source bytes into 4 bytes
	val = (uint32(byte(src[(6+0)/8]>>((6+0)%8*8)))<<16 |
		uint32(byte(src[(6+1)/8]>>((6+1)%8*8)))<<8 |
		uint32(byte(src[(6+2)/8]>>((6+2)%8*8)))<<0)

	dst[8+0] = encodeMap[val>>18&0x3F]
	dst[8+1] = encodeMap[val>>12&0x3F]
	dst[8+2] = encodeMap[val>>6&0x3F]
	dst[8+3] = encodeMap[val&0x3F]

	// Convert 3x 8bit source bytes into 4 bytes
	val = (uint32(byte(src[(9+0)/8]>>((9+0)%8*8)))<<16 |
		uint32(byte(src[(9+1)/8]>>((9+1)%8*8)))<<8 |
		uint32(byte(src[(9+2)/8]>>((9+2)%8*8)))<<0)

	dst[12+0] = encodeMap[val>>18&0x3F]
	dst[12+1] = encodeMap[val>>12&0x3F]
	dst[12+2] = encodeMap[val>>6&0x3F]
	dst[12+3] = encodeMap[val&0x3F]

	// Convert 3x 8bit source bytes into 4 bytes
	val = (uint32(byte(src[(12+0)/8]>>((12+0)%8*8)))<<16 |
		uint32(byte(src[(12+1)/8]>>((12+1)%8*8)))<<8 |
		uint32(byte(src[(12+2)/8]>>((12+2)%8*8)))<<0)

	dst[16+0] = encodeMap[val>>18&0x3F]
	dst[16+1] = encodeMap[val>>12&0x3F]
	dst[16+2] = encodeMap[val>>6&0x3F]
	dst[16+3] = encodeMap[val&0x3F]

	// last byte
	val = uint32(src[1]>>56) << 16
	dst[20+0] = encodeMap[val>>18&0x3F]
	dst[20+1] = encodeMap[val>>12&0x3F]
}
