# randid

```go
import "go.withmatt.com/randid"

id := randid.New().String()
```

`randid`'s goals are simple: a smaller, more compact, faster version of uuid4. I don't care about the actual format of a UUID, I just want some random id.

### comparison with uuid4

* randid is full 128 bits, and uuid4 is 122
* randid is base64 encoded, uuid4 is hex with hyphens
* randid string length is 22, and uuid4 is 32

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
BenchmarkNew
BenchmarkNew-8                   1462352               827.4 ns/op            16 B/op          1 allocs/op
BenchmarkUUIDNew
BenchmarkUUIDNew-8               1458153               833.9 ns/op            16 B/op          1 allocs/op
BenchmarkNewString
BenchmarkNewString-8             1418264               842.0 ns/op            16 B/op          1 allocs/op
BenchmarkUUIDNewString
BenchmarkUUIDNewString-8         1370332               872.0 ns/op            64 B/op          2 allocs/op
```

As we can see, just the `New()` path in comparison to uuid4 is marginally faster, since we're not throwing away extra random bits.

On the `New().String()` path, we gain a bit more performance and 1 less memory allocation, as well as a smaller allocation needed.
