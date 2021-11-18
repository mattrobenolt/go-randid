package randid

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/google/uuid"
	"testing"
)

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New()
	}
	b.ReportAllocs()
}

func BenchmarkUUIDNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = uuid.New()
	}
	b.ReportAllocs()

}

func BenchmarkNewString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New().String()
	}
	b.ReportAllocs()
}

func BenchmarkUUIDNewString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		uuid.New().String()
	}
	b.ReportAllocs()
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

func fakeReader() func([]byte) (int, error) {
	buf := make([]byte, Size)
	rand.Read(buf)

	return func(p []byte) (n int, err error) {
		for i := 0; i < len(p); i++ {
			p[i] = buf[i]
		}
		return len(p), nil
	}
}
