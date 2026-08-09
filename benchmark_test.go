package orderedmap_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/winebarrel/orderedmap/v2"
)

func BenchmarkSet(b *testing.B) {
	om := orderedmap.New[int, struct{}]()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		om.Set(i, struct{}{})
	}
}

func BenchmarkGet(b *testing.B) {
	om := orderedmap.New[int, struct{}]()

	for i := 0; i < b.N; i++ {
		om.Set(i, struct{}{})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		om.Get(i)
	}
}

func BenchmarkDelete(b *testing.B) {
	om := orderedmap.New[int, struct{}]()

	for i := 0; i < b.N; i++ {
		om.Set(i, struct{}{})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		om.Delete(i)
	}
}

func BenchmarkMarshalJSON(b *testing.B) {
	for _, n := range []int{10, 100, 1000} {
		b.Run(fmt.Sprintf("%d keys", n), func(b *testing.B) {
			om := orderedmap.New[string, any]()
			for i := 0; i < n; i++ {
				om.Set(fmt.Sprintf("key%d", i), i)
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				if _, err := json.Marshal(om); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkUnmarshalJSON(b *testing.B) {
	for _, n := range []int{10, 100, 1000} {
		b.Run(fmt.Sprintf("%d keys", n), func(b *testing.B) {
			// Build JSON input once
			buf := []byte{'{'}
			for i := 0; i < n; i++ {
				if i > 0 {
					buf = append(buf, ',')
				}
				buf = append(buf, []byte(fmt.Sprintf(`"key%d":%d`, i, i))...)
			}
			buf = append(buf, '}')

			// Create the ordered map once and reuse it; UnmarshalJSON clears it each time.
			om := orderedmap.New[string, any]()

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				if err := json.Unmarshal(buf, om); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
