package orderedmap

import (
	"container/list"
	"fmt"
	"iter"
	"strings"
	"sync"
)

type OrderedMap[K comparable, V any] struct {
	mu           sync.RWMutex
	pairs        *list.List
	elementByKey map[K]*list.Element
}

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

func New[K comparable, V any]() *OrderedMap[K, V] {
	list.New()

	m := &OrderedMap[K, V]{
		pairs:        list.New(),
		elementByKey: map[K]*list.Element{},
	}

	return m
}

func (om *OrderedMap[K, V]) Len() int {
	om.mu.RLock()
	defer om.mu.RUnlock()

	return len(om.elementByKey)
}

func (om *OrderedMap[K, V]) Set(k K, v V) {
	om.mu.Lock()
	defer om.mu.Unlock()
	p := &Pair[K, V]{Key: k, Value: v}

	if e, ok := om.elementByKey[k]; ok {
		e.Value = p
	} else {
		e = om.pairs.PushBack(p)
		om.elementByKey[k] = e
	}
}

func (om *OrderedMap[K, V]) Push(k K, v V) {
	om.mu.Lock()
	defer om.mu.Unlock()
	p := &Pair[K, V]{Key: k, Value: v}

	if e, ok := om.elementByKey[k]; ok {
		om.pairs.Remove(e)
	}

	e := om.pairs.PushBack(p)
	om.elementByKey[k] = e
}

func (om *OrderedMap[K, V]) Get(k K) V {
	v, _ := om.GetOk(k)
	return v
}

func (om *OrderedMap[K, V]) GetOk(k K) (V, bool) {
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

func (om *OrderedMap[K, V]) Delete(k K) V {
	v, _ := om.DeleteOk(k)
	return v
}

func (om *OrderedMap[K, V]) DeleteOk(k K) (V, bool) {
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

func (om *OrderedMap[K, V]) Clear() {
	om.mu.Lock()
	defer om.mu.Unlock()
	clear(om.elementByKey)
	om.pairs.Init()
}

func (om *OrderedMap[K, V]) Pairs() []*Pair[K, V] {
	om.mu.RLock()
	defer om.mu.RUnlock()

	pairs := []*Pair[K, V]{}
	for e := om.pairs.Front(); e != nil; e = e.Next() {
		p := e.Value.(*Pair[K, V])
		pairs = append(pairs, p)
	}
	return pairs
}

func (om *OrderedMap[K, V]) All() iter.Seq2[K, V] {
	om.mu.RLock()
	defer om.mu.RUnlock()

	return func(yield func(K, V) bool) {
		for e := om.pairs.Front(); e != nil; e = e.Next() {
			p := e.Value.(*Pair[K, V])
			if !yield(p.Key, p.Value) {
				return
			}
		}
	}
}

func (om *OrderedMap[K, V]) Keys() iter.Seq[K] {
	om.mu.RLock()
	defer om.mu.RUnlock()

	return func(yield func(K) bool) {
		for e := om.pairs.Front(); e != nil; e = e.Next() {
			p := e.Value.(*Pair[K, V])
			if !yield(p.Key) {
				return
			}
		}
	}
}

func (om *OrderedMap[K, V]) Values() iter.Seq[V] {
	om.mu.RLock()
	defer om.mu.RUnlock()

	return func(yield func(V) bool) {
		for e := om.pairs.Front(); e != nil; e = e.Next() {
			p := e.Value.(*Pair[K, V])
			if !yield(p.Value) {
				return
			}
		}
	}
}

func (om *OrderedMap[K, V]) String() string {
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
