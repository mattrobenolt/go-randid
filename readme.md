# randid

```go
import "go.withmatt.com/randid"

id := randid.New().String()
```

`randid`'s goals are simple: a smaller, more compact, faster version of uuid4.
I don't care about the actual format of a UUID or being cryptographically secure,
I just want some random id and for it to be fast. The primary intended use cases
are for things like generating a request id, or a transaction id to be used within
logging or metrics. A collision in some wild scenario wouldn't mean a security issue,
just might be a weird coincidence and confusing.

### comparison with uuid4

* randid is full 128 bits, and uuid4 is 122
* randid is base64 encoded, uuid4 is hex with hyphens
* randid string length is 22, and uuid4 is 32
* randid is not cryptographically secure

```
 uuid4: 3cf6702a-da2a-4456-8a53-80a235b3cbfd
randid: cKTij4eSRWmIFgydqgi_Ww
```

### benchmarks

```
goos: linux
goarch: amd64
pkg: go.withmatt.com/randid
cpu: 11th Gen Intel(R) Core(TM) i7-1165G7 @ 2.80GHz
BenchmarkNew-8                       63667344          18.72 ns/op        0 B/op       0 allocs/op
BenchmarkNewString-8                 31164199          39.19 ns/op        0 B/op       0 allocs/op
BenchmarkUUIDNew/NoPool-8             3921735         301.0 ns/op        16 B/op       1 allocs/op
BenchmarkUUIDNew/Pool-8              18226960          68.79 ns/op        0 B/op       0 allocs/op
BenchmarkUUIDNewString/NoPool-8       3260712         373.1 ns/op        64 B/op       2 allocs/op
BenchmarkUUIDNewString/Pool-8        11458695         103.7 ns/op        48 B/op       1 allocs/op

```
In general, our performance is significantly faster on all counts fundamentally by not using
cryptographically secure randoms. Random bytes are generated as efficiently as possible, and focus
is on avoiding allocations and unnecessary instructions.

On the `New().String()` path, we also avoid an allocation and is significantly faster than
simply base64 encoding.
