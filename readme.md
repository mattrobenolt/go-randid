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
goos: darwin
goarch: arm64
pkg: go.withmatt.com/randid
BenchmarkNew
BenchmarkNew-10            100000000            10.06 ns/op        0 B/op      0 allocs/op
BenchmarkString
BenchmarkString-10         121892878            9.857 ns/op        0 B/op      0 allocs/op
BenchmarkNewString
BenchmarkNewString-10       48996674            24.39 ns/op        0 B/op      0 allocs/op
```

```
$ benchstat uuid.txt randid.txt
             │   uuid.txt   │             randid.txt              │
             │    sec/op    │   sec/op     vs base                │
New-10          46.99n ± 0%   10.05n ± 1%  -78.60% (p=0.000 n=10)
String-10      34.435n ± 6%   9.839n ± 0%  -71.43% (p=0.000 n=10)
NewString-10    75.98n ± 0%   24.56n ± 0%  -67.68% (p=0.000 n=10)
geomean         49.72n        13.44n       -72.96%

             │   uuid.txt   │               randid.txt                │
             │     B/op     │    B/op     vs base                     │
New-10         0.000 ± 0%     0.000 ± 0%         ~ (p=1.000 n=10) ¹
String-10      48.00 ± 0%      0.00 ± 0%  -100.00% (p=0.000 n=10)
NewString-10   48.00 ± 0%      0.00 ± 0%  -100.00% (p=0.000 n=10)
geomean                   ²               ?                       ² ³
¹ all samples are equal
² summaries must be >0 to compute geomean
³ ratios must be >0 to compute geomean

             │   uuid.txt   │               randid.txt                │
             │  allocs/op   │ allocs/op   vs base                     │
New-10         0.000 ± 0%     0.000 ± 0%         ~ (p=1.000 n=10) ¹
String-10      1.000 ± 0%     0.000 ± 0%  -100.00% (p=0.000 n=10)
NewString-10   1.000 ± 0%     0.000 ± 0%  -100.00% (p=0.000 n=10)
geomean                   ²               ?                       ² ³
¹ all samples are equal
² summaries must be >0 to compute geomean
³ ratios must be >0 to compute geomean
```

In general, our performance is significantly faster on all counts fundamentally by not using
cryptographically secure randoms. Random bytes are generated as efficiently as possible, and focus
is on avoiding allocations and unnecessary instructions.

On the `New().String()` path, we also avoid an allocation and is significantly faster than
simply base64 encoding.
