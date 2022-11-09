package randid

import (
	"crypto/rand"
	"encoding/base64"
	"testing"

	"github.com/google/uuid"
)

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New()
	}
	b.ReportAllocs()
}

func BenchmarkUUIDNew(b *testing.B) {
	b.Run("NoPool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = uuid.New()
		}
		b.ReportAllocs()
	})
	b.Run("Pool", func(b *testing.B) {
		uuid.EnableRandPool()
		defer uuid.DisableRandPool()
		for i := 0; i < b.N; i++ {
			_ = uuid.New()
		}
		b.ReportAllocs()
	})
}

func BenchmarkNewString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New().String()
	}
	b.ReportAllocs()
}

func BenchmarkUUIDNewString(b *testing.B) {
	b.Run("NoPool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = uuid.New().String()
		}
		b.ReportAllocs()
	})
	b.Run("Pool", func(b *testing.B) {
		uuid.EnableRandPool()
		defer uuid.DisableRandPool()
		for i := 0; i < b.N; i++ {
			_ = uuid.New().String()
		}
		b.ReportAllocs()
	})
}

func TestEncode(t *testing.T) {
	currentReader := randReader
	defer func() { randReader = currentReader }()
	randReader = fakeReader()

	buf := make([]byte, Size)
	randReader(buf)
	expected := base64.RawURLEncoding.EncodeToString(buf)

	got := New().String()

	if len(got) != StringLen {
		t.Errorf("Expected length %d, got %d", StringLen, len(got))
	}

	if got != expected {
		t.Errorf("\nexpected: %s\n     got: %s", expected, got)
	}
}

func TestNew(t *testing.T) {
	currentReader := randReader
	defer func() { randReader = currentReader }()
	recorder := recordedReader(randReader)
	randReader = recorder.read

	const expectedReads = 3
	const ids = idsPerPage * expectedReads

	results := make(map[string]struct{}, ids)
	for i := 0; i < ids; i++ {
		id := New().String()
		if _, ok := results[id]; ok {
			t.Errorf("Duplicate id generated")
		}
		results[id] = struct{}{}
	}

	if len(results) != ids {
		t.Errorf("Expected ids %d, got %d", idsPerPage*expectedReads, len(results))
	}

	if recorder.count != expectedReads {
		t.Errorf("Expected reads %d, got %d", expectedReads, recorder.count)
	}
}

func fakeReader() func([]byte) (int, error) {
	buf := make([]byte, pageSize)
	rand.Read(buf)

	return func(p []byte) (n int, err error) {
		for i := 0; i < len(p); i++ {
			p[i] = buf[i]
		}
		return len(p), nil
	}
}

type readRecorder struct {
	count int
	fn    func([]byte) (int, error)
}

func (r *readRecorder) read(p []byte) (int, error) {
	r.count++
	return r.fn(p)
}

func recordedReader(fn func([]byte) (int, error)) *readRecorder {
	return &readRecorder{fn: fn}
}
