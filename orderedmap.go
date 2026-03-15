package orderedmap

import (
	"container/list"
	"fmt"
	"iter"
	"slices"
	"strings"
	"sync"
)

type Map[K comparable, V any] struct {
	mu           sync.RWMutex
	entries      *list.List
	elementByKey map[K]*list.Element
}

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

func New[K comparable, V any]() *Map[K, V] {
	om := &Map[K, V]{
		entries:      list.New(),
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

func (om *Map[K, V]) Clone() *Map[K, V] {
	return From(om.All())
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
		e = om.entries.PushBack(p)
		om.elementByKey[k] = e
	}
}

func (om *Map[K, V]) Push(k K, v V) {
	om.mu.Lock()
	defer om.mu.Unlock()
	p := &Pair[K, V]{Key: k, Value: v}

	if e, ok := om.elementByKey[k]; ok {
		om.entries.Remove(e)
	}

	e := om.entries.PushBack(p)
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
		om.entries.Remove(e)
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
	om.entries.Init()
}

func (om *Map[K, V]) Entries() []Pair[K, V] {
	om.mu.RLock()
	defer om.mu.RUnlock()

	return om.entries0()
}

func (om *Map[K, V]) entries0() []Pair[K, V] {
	pairs := make([]Pair[K, V], 0, om.entries.Len())
	for e := om.entries.Front(); e != nil; e = e.Next() {
		p := e.Value.(*Pair[K, V])
		pairs = append(pairs, *p)
	}
	return pairs
}

func (om *Map[K, V]) All() iter.Seq2[K, V] {
	pairs := om.Entries()

	return func(yield func(K, V) bool) {
		for _, p := range pairs {
			if !yield(p.Key, p.Value) {
				return
			}
		}
	}
}

func (om *Map[K, V]) Keys() iter.Seq[K] {
	pairs := om.Entries()

	return func(yield func(K) bool) {
		for _, p := range pairs {
			if !yield(p.Key) {
				return
			}
		}
	}
}

func (om *Map[K, V]) CollectKeys() []K {
	if om.Len() == 0 {
		return []K{}
	} else {
		return slices.Collect(om.Keys())
	}
}

func (om *Map[K, V]) Values() iter.Seq[V] {
	pairs := om.Entries()

	return func(yield func(V) bool) {
		for _, p := range pairs {
			if !yield(p.Value) {
				return
			}
		}
	}
}

func (om *Map[K, V]) CollectValues() []V {
	if om.Len() == 0 {
		return []V{}
	} else {
		return slices.Collect(om.Values())
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
