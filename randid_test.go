package randid

import (
	"testing"
)

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New()
	}
	b.ReportAllocs()
}

func BenchmarkString(b *testing.B) {
	id := New()
	for i := 0; i < b.N; i++ {
		_ = id.String()
	}
	b.ReportAllocs()
}

func BenchmarkNewString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New().String()
	}
	b.ReportAllocs()
}

// func BenchmarkNewParallel(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		var wg sync.WaitGroup
// 		for range 10 {
// 			wg.Add(1)
// 			go func() {
// 				_ = New()
// 				wg.Done()
// 			}()
// 		}
// 		wg.Wait()
// 	}
// 	b.ReportAllocs()
// }

// func BenchmarkNew(b *testing.B) {
// 	uuid.EnableRandPool()
// 	defer uuid.DisableRandPool()
// 	for i := 0; i < b.N; i++ {
// 		_ = uuid.New()
// 	}
// 	b.ReportAllocs()
// }

// func BenchmarkString(b *testing.B) {
// 	uuid.EnableRandPool()
// 	defer uuid.DisableRandPool()

// 	id := uuid.New()
// 	for i := 0; i < b.N; i++ {
// 		_ = id.String()
// 	}
// 	b.ReportAllocs()
// }

// func BenchmarkNewString(b *testing.B) {
// 	uuid.EnableRandPool()
// 	defer uuid.DisableRandPool()
// 	for i := 0; i < b.N; i++ {
// 		_ = uuid.New().String()
// 	}
// 	b.ReportAllocs()
// }

func TestNew(t *testing.T) {
	const ids = 100

	results := make(map[ID]struct{}, ids)
	for i := 0; i < ids; i++ {
		id := New()
		if _, ok := results[id]; ok {
			t.Errorf("Duplicate id generated")
		}
		results[id] = struct{}{}
	}

	if len(results) != ids {
		t.Errorf("Expected ids %d, got %d", ids, len(results))
	}
}

func TestEncode(t *testing.T) {
	cases := []struct {
		id       ID
		expected string
	}{
		{ID{0, 0}, "AAAAAAAAAAAAAAAAAAAAAA"},
		{ID{1, 0}, "AQAAAAAAAAAAAAAAAAAAAA"},
		{ID{0, 1}, "AAAAAAAAAAABAAAAAAAAAA"},
		{ID{1485390858579739336, 6272678278200745536}, "yAJCrh0rnRRA3nbyqAcNVw"},
		{ID{18446744073709551615, 18446744073709551615}, "_____________________w"},
	}

	for _, c := range cases {
		got := c.id.String()
		if len(got) != StringLen {
			t.Errorf("Expected length %d, got %d", StringLen, len(got))
		}
		if got != c.expected {
			t.Errorf("Expected %v, got %v (%#v)", c.expected, got, c.id)
		}
	}
}
