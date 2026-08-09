# orderedmap

[![GitHub Tag](https://img.shields.io/github/v/tag/winebarrel/orderedmap.svg)](https://pkg.go.dev/github.com/winebarrel/orderedmap?tab=versions)
[![Go Reference](https://pkg.go.dev/badge/github.com/winebarrel/orderedmap/v2.svg)](https://pkg.go.dev/github.com/winebarrel/orderedmap/v2)
[![CI](https://github.com/winebarrel/orderedmap/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/orderedmap/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/winebarrel/orderedmap/graph/badge.svg)](https://codecov.io/gh/winebarrel/orderedmap)

A generic, thread-safe ordered map for Go that preserves insertion order using [linkedlist](https://github.com/winebarrel/linkedlist).

## Installation

```
go get github.com/winebarrel/orderedmap/v2
```

Requires Go 1.25+.

## Usage

```go
package main

import (
	"fmt"
	"slices"

	"github.com/winebarrel/orderedmap/v2"
)

func main() {
	om := orderedmap.New[string, any]()

	om.Set("foo", "bar")
	om.Set("zoo", 100)
	om.Set("baz", true)

	fmt.Println(om)
	//=> *orderedmap.Map[string,interface {}][foo:bar zoo:100 baz:true]

	// Get a value
	fmt.Println(om.Get("foo"))
	//=> bar

	// Get a value with existence check
	v, ok := om.GetOk("foo")
	fmt.Println(v, ok)
	//=> bar true

	// Iterate over key-value pairs in insertion order
	for k, v := range om.All() {
		fmt.Println(k, v)
	}
	//=> foo bar
	//   zoo 100
	//   baz true

	// Iterate over keys
	for k := range om.Keys() {
		fmt.Println(k)
	}

	// Iterate over values
	for v := range om.Values() {
		fmt.Println(v)
	}

	// Push moves an existing key to the back
	om.Push("foo", "new_bar")
	for k, v := range om.All() {
		fmt.Println(k, v)
	}
	//=> zoo 100
	//   baz true
	//   foo new_bar

	// Get all key-value pairs as a slice
	pairs := om.Entries()
	fmt.Println(pairs)
	//=> [{zoo 100} {baz true} {foo new_bar}]

	// Delete a key
	om.Delete("zoo")

	// Get the number of entries
	fmt.Println(om.Len())
	//=> 2

	// Clear all entries
	om.Clear()

	// Create from an iterator
	om2 := orderedmap.From(slices.All([]string{"foo", "bar", "baz"}))
	fmt.Println(om2)
	//=> *orderedmap.Map[int,string][0:foo 1:bar 2:baz]
}
```

## Transform

```go
package main

import (
	"fmt"

	"github.com/winebarrel/orderedmap/v2"
)

func main() {
	om := orderedmap.New[string, int]()
	om.Set("foo", 1)
	om.Set("bar", 2)
	om.Set("baz", 3)

	// Transform returns an iterator of transformed values
	for s := range orderedmap.Transform(om, func(k string, v int) string {
		return fmt.Sprintf("%s=%d", k, v)
	}) {
		fmt.Println(s)
	}
	//=> foo=1
	//   bar=2
	//   baz=3

	// TransformSlice returns a slice of transformed values
	ss := orderedmap.TransformSlice(om, func(k string, v int) string {
		return fmt.Sprintf("%s=%d", k, v)
	})
	fmt.Println(ss)
	//=> [foo=1 bar=2 baz=3]
}
```

## JSON Marshal / Unmarshal

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/winebarrel/orderedmap/v2"
)

func main() {
	// Marshal: preserves insertion order
	om := orderedmap.New[string, any]()
	om.Set("z", 3)
	om.Set("a", 1)
	om.Set("m", 2)

	b, _ := json.Marshal(om)
	fmt.Println(string(b))
	//=> {"z":3,"a":1,"m":2}

	// Unmarshal: preserves key order from JSON
	om2 := orderedmap.New[string, any]()
	json.Unmarshal([]byte(`{"z":3,"a":1,"m":2}`), om2)

	for k, v := range om2.All() {
		fmt.Println(k, v)
	}
	//=> z 3
	//   a 1
	//   m 2
}
```
