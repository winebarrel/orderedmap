package orderedmap

import (
	"cmp"
	"container/list"
	"fmt"
	"iter"
	"slices"
	"strings"
	"sync"
)

type Map[K comparable, V any] struct {
	mu           sync.RWMutex
	pairs        *list.List
	elementByKey map[K]*list.Element
}

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

func New[K comparable, V any]() *Map[K, V] {
	om := &Map[K, V]{
		pairs:        list.New(),
		elementByKey: map[K]*list.Element{},
	}

	return om
}

func From[K comparable, V any](seq iter.Seq2[K, V]) *Map[K, V] {
	om := New[K, V]()

	for k, v := range seq {
		om.set0(k, v)
	}

	return om
}

func (om *Map[K, V]) Len() int {
	om.mu.RLock()
	defer om.mu.RUnlock()

	return len(om.elementByKey)
}

func (om *Map[K, V]) Set(k K, v V) {
	om.mu.Lock()
	defer om.mu.Unlock()

	om.set0(k, v)
}

func (om *Map[K, V]) set0(k K, v V) {
	p := &Pair[K, V]{Key: k, Value: v}

	if e, ok := om.elementByKey[k]; ok {
		e.Value = p
	} else {
		e = om.pairs.PushBack(p)
		om.elementByKey[k] = e
	}
}

func (om *Map[K, V]) Push(k K, v V) {
	om.mu.Lock()
	defer om.mu.Unlock()
	p := &Pair[K, V]{Key: k, Value: v}

	if e, ok := om.elementByKey[k]; ok {
		om.pairs.Remove(e)
	}

	e := om.pairs.PushBack(p)
	om.elementByKey[k] = e
}

func (om *Map[K, V]) Get(k K) V {
	v, _ := om.GetOk(k)
	return v
}

func (om *Map[K, V]) GetOk(k K) (V, bool) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	if e, ok := om.elementByKey[k]; ok {
		p := e.Value.(*Pair[K, V])
		return p.Value, true
	} else {
		var v V
		return v, false
	}
}

func (om *Map[K, V]) Delete(k K) V {
	v, _ := om.DeleteOk(k)
	return v
}

func (om *Map[K, V]) DeleteOk(k K) (V, bool) {
	om.mu.Lock()
	defer om.mu.Unlock()

	if e, ok := om.elementByKey[k]; ok {
		delete(om.elementByKey, k)
		om.pairs.Remove(e)
		p := e.Value.(*Pair[K, V])
		return p.Value, true
	} else {
		var v V
		return v, false
	}
}

func (om *Map[K, V]) Clear() {
	om.mu.Lock()
	defer om.mu.Unlock()
	clear(om.elementByKey)
	om.pairs.Init()
}

func (om *Map[K, V]) Pairs() []Pair[K, V] {
	om.mu.RLock()
	defer om.mu.RUnlock()

	pairs := make([]Pair[K, V], 0, om.pairs.Len())
	for e := om.pairs.Front(); e != nil; e = e.Next() {
		p := e.Value.(*Pair[K, V])
		pairs = append(pairs, *p)
	}
	return pairs
}

func (om *Map[K, V]) All() iter.Seq2[K, V] {
	pairs := om.Pairs()

	return func(yield func(K, V) bool) {
		for _, p := range pairs {
			if !yield(p.Key, p.Value) {
				return
			}
		}
	}
}

func (om *Map[K, V]) Keys() iter.Seq[K] {
	pairs := om.Pairs()

	return func(yield func(K) bool) {
		for _, p := range pairs {
			if !yield(p.Key) {
				return
			}
		}
	}
}

func (om *Map[K, V]) Values() iter.Seq[V] {
	pairs := om.Pairs()

	return func(yield func(V) bool) {
		for _, p := range pairs {
			if !yield(p.Value) {
				return
			}
		}
	}
}

func (om *Map[K, V]) String() string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "%T[", om)
	first := true
	for k, v := range om.All() {
		if !first {
			buf.WriteString(" ")
		}
		fmt.Fprintf(&buf, "%v:%v", k, v)
		first = false
	}
	buf.WriteString("]")

	return buf.String()
}

func Sorted[K cmp.Ordered, V any](om *Map[K, V]) *Map[K, V] {
	return SortedFunc(om, func(a, b Pair[K, V]) int {
		return cmp.Compare(a.Key, b.Key)
	})
}

func SortedFunc[K cmp.Ordered, V any](om *Map[K, V], cmp func(Pair[K, V], Pair[K, V]) int) *Map[K, V] {
	pairs := om.Pairs()
	slices.SortFunc(pairs, cmp)

	return From(func(yield func(K, V) bool) {
		for _, p := range pairs {
			if !yield(p.Key, p.Value) {
				return
			}
		}
	})
}
