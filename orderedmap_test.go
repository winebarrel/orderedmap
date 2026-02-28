package orderedmap_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/orderedmap"
)

type pair[V any] struct {
	k string
	v V
}

type expectedOk struct {
	v  any
	ok bool
}

func mapToPairs[V any](t *testing.T, om *orderedmap.OrderedMap[string, V]) []pair[V] {
	t.Helper()
	pairs := []pair[V]{}
	for k, v := range om.All() {
		pairs = append(pairs, pair[V]{k: k, v: v})
	}
	return pairs
}

func TestLen(t *testing.T) {
	tests := []struct {
		name     string
		init     []pair[any]
		expected int
	}{
		{
			name:     "set value",
			init:     []pair[any]{},
			expected: 0,
		},
		{
			name:     "set value",
			init:     []pair[any]{{k: "foo", v: "bar"}},
			expected: 1,
		},
		{
			name:     "set value",
			init:     []pair[any]{{k: "foo", v: "bar"}, {k: "zoo", v: 100}},
			expected: 2,
		},
		{
			name:     "set value",
			init:     []pair[any]{{k: "foo", v: "bar"}, {k: "zoo", v: 100}, {k: "baz", v: true}},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := orderedmap.New[string, any]()
			for _, p := range tt.init {
				om.Set(p.k, p.v)
			}
			assert.Equal(t, tt.expected, om.Len())
		})
	}
}

func TestSet(t *testing.T) {
	tests := []struct {
		name     string
		init     []pair[any]
		set      []pair[any]
		expected []pair[any]
	}{
		{
			name:     "set value",
			init:     []pair[any]{},
			set:      []pair[any]{{k: "foo", v: "bar"}, {k: "zoo", v: 100}, {k: "baz", v: true}},
			expected: []pair[any]{{k: "foo", v: "bar"}, {k: "zoo", v: 100}, {k: "baz", v: true}},
		},
		{
			name:     "set exists value",
			init:     []pair[any]{{k: "zoo", v: 100}, {k: "baz", v: true}},
			set:      []pair[any]{{k: "foo", v: "bar"}, {k: "zoo", v: 200}, {k: "baz", v: false}, {k: "hoge", v: "fuga"}},
			expected: []pair[any]{{k: "zoo", v: 200}, {k: "baz", v: false}, {k: "foo", v: "bar"}, {k: "hoge", v: "fuga"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := orderedmap.New[string, any]()
			for _, p := range tt.init {
				om.Set(p.k, p.v)
			}
			for _, p := range tt.set {
				om.Set(p.k, p.v)
			}
			for range 10 {
				assert.Equal(t, tt.expected, mapToPairs(t, om))
			}
		})
	}
}

func TestPush(t *testing.T) {
	tests := []struct {
		name     string
		init     []pair[any]
		push     []pair[any]
		expected []pair[any]
	}{
		{
			name:     "push value",
			init:     []pair[any]{},
			push:     []pair[any]{{k: "foo", v: "bar"}, {k: "zoo", v: 100}, {k: "baz", v: true}},
			expected: []pair[any]{{k: "foo", v: "bar"}, {k: "zoo", v: 100}, {k: "baz", v: true}},
		},
		{
			name:     "set exists value",
			init:     []pair[any]{{k: "zoo", v: 100}, {k: "baz", v: true}},
			push:     []pair[any]{{k: "foo", v: "bar"}, {k: "zoo", v: 200}, {k: "baz", v: false}, {k: "hoge", v: "fuga"}},
			expected: []pair[any]{{k: "foo", v: "bar"}, {k: "zoo", v: 200}, {k: "baz", v: false}, {k: "hoge", v: "fuga"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := orderedmap.New[string, any]()
			for _, p := range tt.init {
				om.Set(p.k, p.v)
			}
			for _, p := range tt.push {
				om.Push(p.k, p.v)
			}
			for range 10 {
				assert.Equal(t, tt.expected, mapToPairs(t, om))
			}
		})
	}
}

func TestGetAny(t *testing.T) {
	tests := []struct {
		name           string
		init           []pair[any]
		get            []string
		expectedValues []any
	}{
		{
			name:           "get value",
			init:           []pair[any]{{k: "foo", v: "bar"}, {k: "baz", v: true}},
			get:            []string{"foo", "bar", "baz"},
			expectedValues: []any{"bar", (any)(nil), true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := orderedmap.New[string, any]()
			for _, p := range tt.init {
				om.Set(p.k, p.v)
			}
			for range 10 {
				values := []any{}
				for _, k := range tt.get {
					values = append(values, om.Get(k))
				}
				assert.Equal(t, tt.expectedValues, values)
				assert.Equal(t, tt.init, mapToPairs(t, om))
			}
		})
	}
}

func TestGetInt(t *testing.T) {
	tests := []struct {
		name           string
		init           []pair[int]
		get            []string
		expectedValues []any
	}{
		{
			name:           "get value",
			init:           []pair[int]{{k: "foo", v: 100}, {k: "baz", v: 200}},
			get:            []string{"foo", "bar", "baz"},
			expectedValues: []any{100, 0, 200},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := orderedmap.New[string, int]()
			for _, p := range tt.init {
				om.Set(p.k, p.v)
			}
			for range 10 {
				values := []any{}
				for _, k := range tt.get {
					values = append(values, om.Get(k))
				}
				assert.Equal(t, tt.expectedValues, values)
				assert.Equal(t, tt.init, mapToPairs(t, om))
			}
		})
	}
}

func TestGetString(t *testing.T) {
	tests := []struct {
		name           string
		init           []pair[string]
		get            []string
		expectedValues []any
	}{
		{
			name:           "get value",
			init:           []pair[string]{{k: "foo", v: "bar"}, {k: "baz", v: "zoo"}},
			get:            []string{"foo", "bar", "baz"},
			expectedValues: []any{"bar", "", "zoo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := orderedmap.New[string, string]()
			for _, p := range tt.init {
				om.Set(p.k, p.v)
			}
			for range 10 {
				values := []any{}
				for _, k := range tt.get {
					values = append(values, om.Get(k))
				}
				assert.Equal(t, tt.expectedValues, values)
				assert.Equal(t, tt.init, mapToPairs(t, om))
			}
		})
	}
}

func TestGetOK(t *testing.T) {
	tests := []struct {
		name           string
		init           []pair[any]
		get            []string
		expectedValues []expectedOk
	}{
		{
			name:           "get value",
			init:           []pair[any]{{k: "foo", v: "bar"}, {k: "baz", v: true}},
			get:            []string{"foo", "bar", "baz"},
			expectedValues: []expectedOk{{v: "bar", ok: true}, {v: (any)(nil), ok: false}, {v: true, ok: true}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := orderedmap.New[string, any]()
			for _, p := range tt.init {
				om.Set(p.k, p.v)
			}
			for range 10 {
				es := []expectedOk{}
				for _, k := range tt.get {
					e := expectedOk{}
					e.v, e.ok = om.GetOk(k)
					es = append(es, e)
				}
				assert.Equal(t, tt.expectedValues, es)
				assert.Equal(t, tt.init, mapToPairs(t, om))
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name           string
		init           []pair[any]
		delete         []string
		expectedMap    []pair[any]
		expectedValues []any
	}{
		{
			name:           "delete value",
			init:           []pair[any]{{k: "foo", v: "bar"}, {k: "baz", v: true}},
			delete:         []string{"foo", "bar"},
			expectedMap:    []pair[any]{{k: "baz", v: true}},
			expectedValues: []any{"bar", (any)(nil)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := orderedmap.New[string, any]()
			for _, p := range tt.init {
				om.Set(p.k, p.v)
			}
			values := []any{}
			for _, k := range tt.delete {
				values = append(values, om.Delete(k))
			}
			assert.Equal(t, tt.expectedMap, mapToPairs(t, om))
			assert.Equal(t, tt.expectedValues, values)
		})
	}
}

func TestDeleteOK(t *testing.T) {
	tests := []struct {
		name           string
		init           []pair[any]
		delete         []string
		expectedMap    []pair[any]
		expectedValues []expectedOk
	}{
		{
			name:           "delete value",
			init:           []pair[any]{{k: "foo", v: "bar"}, {k: "baz", v: true}},
			delete:         []string{"foo", "bar"},
			expectedMap:    []pair[any]{{k: "baz", v: true}},
			expectedValues: []expectedOk{{v: "bar", ok: true}, {v: (any)(nil), ok: false}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := orderedmap.New[string, any]()
			for _, p := range tt.init {
				om.Set(p.k, p.v)
			}
			es := []expectedOk{}
			for _, k := range tt.delete {
				e := expectedOk{}
				e.v, e.ok = om.DeleteOk(k)
				es = append(es, e)
			}
			assert.Equal(t, tt.expectedMap, mapToPairs(t, om))
			assert.Equal(t, tt.expectedValues, es)
		})
	}
}

func TestClear(t *testing.T) {
	tests := []struct {
		name        string
		init        []pair[any]
		expectedMap []pair[any]
	}{
		{
			name:        "clear values",
			init:        []pair[any]{{k: "foo", v: "bar"}, {k: "baz", v: true}},
			expectedMap: []pair[any]{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := orderedmap.New[string, any]()
			for _, p := range tt.init {
				om.Set(p.k, p.v)
			}
			om.Clear()
			assert.Equal(t, tt.expectedMap, mapToPairs(t, om))
		})
	}
}

func TestPairs(t *testing.T) {
	tests := []struct {
		name     string
		init     []pair[any]
		expected []pair[any]
	}{
		{
			name:     "all pairs",
			init:     []pair[any]{{k: "foo", v: "bar"}, {k: "zoo", v: 100}, {k: "baz", v: true}},
			expected: []pair[any]{{k: "foo", v: "bar"}, {k: "zoo", v: 100}, {k: "baz", v: true}},
		},
		{
			name:     "all pairs (empty}",
			init:     []pair[any]{},
			expected: []pair[any]{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := orderedmap.New[string, any]()
			for _, p := range tt.init {
				om.Set(p.k, p.v)
			}
			for range 10 {
				pairs := []pair[any]{}
				for _, p := range om.Pairs() {
					pairs = append(pairs, pair[any]{k: p.Key, v: p.Value})
				}
				assert.Equal(t, tt.expected, pairs)
			}
		})
	}
}

func TestKeys(t *testing.T) {
	tests := []struct {
		name     string
		init     []pair[any]
		expected []string
	}{
		{
			name:     "all keys",
			init:     []pair[any]{{k: "foo", v: "bar"}, {k: "zoo", v: 100}, {k: "baz", v: true}},
			expected: []string{"foo", "zoo", "baz"},
		},
		{
			name:     "all keys (empty}",
			init:     []pair[any]{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := orderedmap.New[string, any]()
			for _, p := range tt.init {
				om.Set(p.k, p.v)
			}
			for range 10 {
				keys := []string{}
				for k := range om.Keys() {
					keys = append(keys, k)
				}
				assert.Equal(t, tt.expected, keys)
			}
		})
	}
}

func TestValues(t *testing.T) {
	tests := []struct {
		name     string
		init     []pair[any]
		expected []any
	}{
		{
			name:     "all values",
			init:     []pair[any]{{k: "foo", v: "bar"}, {k: "zoo", v: 100}, {k: "baz", v: true}},
			expected: []any{"bar", 100, true},
		},
		{
			name:     "all values (empty}",
			init:     []pair[any]{},
			expected: []any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := orderedmap.New[string, any]()
			for _, p := range tt.init {
				om.Set(p.k, p.v)
			}
			for range 10 {
				values := []any{}
				for v := range om.Values() {
					values = append(values, v)
				}
				assert.Equal(t, tt.expected, values)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		init     []pair[any]
		expected string
	}{
		{
			name:     "map to string",
			init:     []pair[any]{{k: "foo", v: "bar"}, {k: "zoo", v: 100}, {k: "baz", v: true}},
			expected: `*orderedmap.OrderedMap[string,interface {}][foo:bar zoo:100 baz:true]`,
		},
		{
			name:     "empty map to string",
			init:     []pair[any]{},
			expected: `*orderedmap.OrderedMap[string,interface {}][]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := orderedmap.New[string, any]()
			for _, p := range tt.init {
				om.Set(p.k, p.v)
			}
			for range 10 {
				assert.Equal(t, tt.expected, om.String())
			}
		})
	}
}
