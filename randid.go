package randid

import (
	"math/rand"
	"sync"
	"time"
)

// Size is the length in bytes of the ID
const Size = 16

// StringLen is the length of the string representation of ID
const StringLen = 22

// ID is our 16 byte random value
type ID [Size]byte

// Maintain a pool of rand readers, since these are not safe for
// concurrent use and avoid using a single with a mutex.
var randPool = sync.Pool{
	New: func() interface{} {
		return rand.New(rand.NewSource(time.Now().UnixNano()))
	},
}

// String returns base64 encoding of our ID
func (id ID) String() string {
	var buf [StringLen]byte
	encode(buf[:], id)
	return string(buf[:])
}

var randReader = func() (int64, int64) {
	r := randPool.Get().(*rand.Rand)
	defer randPool.Put(r)
	return r.Int63(), r.Int63()
}

// New generates a new random ID
func New() (id ID) {
	i, j := randReader()

	// reading needs to fill up 16 bytes, or 128 bits of data
	// this means we need to generate 2 x int64s, and put each int64 into 8 bytes

	// hand unrolling since it's not a ton
	id[0] = byte(i)
	id[1] = byte(i >> 8)
	id[2] = byte(i >> 16)
	id[3] = byte(i >> 24)
	id[4] = byte(i >> 32)
	id[5] = byte(i >> 40)
	id[6] = byte(i >> 48)
	id[7] = byte(i >> 56)

	id[8] = byte(j)
	id[9] = byte(j >> 8)
	id[10] = byte(j >> 16)
	id[11] = byte(j >> 24)
	id[12] = byte(j >> 32)
	id[13] = byte(j >> 40)
	id[14] = byte(j >> 48)
	id[15] = byte(j >> 56)
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
