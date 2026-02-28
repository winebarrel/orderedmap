package orderedmap_test

import (
	"testing"

	"github.com/winebarrel/orderedmap"
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
