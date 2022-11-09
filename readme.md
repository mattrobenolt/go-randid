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
BenchmarkNew-8                       19509478          59.18 ns/op        0 B/op       0 allocs/op
BenchmarkUUIDNew/NoPool-8             4156387         291.3 ns/op        16 B/op       1 allocs/op
BenchmarkUUIDNew/Pool-8              18137985          65.17 ns/op        0 B/op       0 allocs/op
BenchmarkNewString-8                 15435547          73.39 ns/op        0 B/op       0 allocs/op
BenchmarkUUIDNewString/NoPool-8       3602911         331.4 ns/op        64 B/op       2 allocs/op
BenchmarkUUIDNewString/Pool-8        11810478          99.93 ns/op       48 B/op       1 allocs/op
```

By default, `New()` uses a buffer pool to optimize reads. The uuid4 library has a toggle to enable or disable their pool for security purposes.

To compare the two, `New()` beats out uuid4 with the pool turned on marginally, and drastically beats using without a pool.

On the `New().String()` path, we gain a bit more performance and 1 less memory allocation, as well as a smaller allocation needed.
