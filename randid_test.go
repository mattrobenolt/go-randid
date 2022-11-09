package randid

import (
	"bytes"
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
	// currentReader := randReader
	// defer func() { randReader = currentReader }()
	// randReader = fakeReader()

	buf := New()
	expected := base64.RawURLEncoding.EncodeToString(buf[:])

	got := buf.String()

	if len(got) != StringLen {
		t.Errorf("Expected length %d, got %d", StringLen, len(got))
	}

	if got != expected {
		t.Errorf("\nexpected: %s\n     got: %s", expected, got)
	}
}

func TestNew(t *testing.T) {
	const ids = 100

	results := make(map[string]struct{}, ids)
	for i := 0; i < ids; i++ {
		id := New().String()
		if _, ok := results[id]; ok {
			t.Errorf("Duplicate id generated")
		}
		results[id] = struct{}{}
	}

	if len(results) != ids {
		t.Errorf("Expected ids %d, got %d", ids, len(results))
	}
}

func TestWithKnownBytes(t *testing.T) {
	currentReader := randReader
	defer func() { randReader = currentReader }()

	cases := []struct {
		i, j     int64
		expected [Size]byte
	}{
		{0, 0, [Size]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{0, 1, [Size]byte{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0}},
		{1, 0, [Size]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{1, 1, [Size]byte{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0}},
		{2 << 8, 4 << 32, [Size]byte{0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0}},
		{1<<0 + (1 << 8), 4 << 32, [Size]byte{1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0}},
		{
			// (1 << 0) + (2 << 8) + (3 << 16) + (4 << 24) + (5 << 32) + (6 << 40) + (7 << 48) + (8 << 56),
			578437695752307201,
			// (9 << 0) + (10 << 8) + (11 << 16) + (12 << 24) + (13 << 32) + (14 << 40) + (15 << 48) + (16 << 56),
			1157159078456920585,
			[Size]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		},
	}

	for _, c := range cases {
		randReader = func() (int64, int64) { return c.i, c.j }
		id := New()
		if !bytes.Equal(id[:], c.expected[:]) {
			t.Errorf("Expected %v, got %v", c.expected, id[:])
		}
	}
}
