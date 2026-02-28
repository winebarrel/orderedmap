package orderedmap_test

import (
	"testing"

	"github.com/winebarrel/orderedmap"
)

func BenchmarkSet(b *testing.B) {
	om := orderedmap.New[int, struct{}]()
	b.ResetTimer()

	for i := range b.N {
		om.Set(i, struct{}{})
	}
}

func BenchmarkGet(b *testing.B) {
	om := orderedmap.New[int, struct{}]()

	for i := range b.N {
		om.Set(i, struct{}{})
	}

	b.ResetTimer()

	for i := range b.N {
		om.Get(i)
	}
}

func BenchmarkDelete(b *testing.B) {
	om := orderedmap.New[int, struct{}]()

	for i := range b.N {
		om.Set(i, struct{}{})
	}

	b.ResetTimer()

	for i := range b.N {
		om.Delete(i)
	}
}
