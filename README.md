# orderedmap

[![GitHub Tag](https://img.shields.io/github/v/tag/winebarrel/orderedmap.svg)](https://pkg.go.dev/github.com/winebarrel/orderedmap?tab=versions)
[![Go Reference](https://pkg.go.dev/badge/github.com/winebarrel/orderedmap.svg)](https://pkg.go.dev/github.com/winebarrel/orderedmap)
[![CI](https://github.com/winebarrel/orderedmap/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/orderedmap/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/winebarrel/orderedmap/graph/badge.svg)](https://codecov.io/gh/winebarrel/orderedmap)
[![Go Report Card](https://goreportcard.com/badge/github.com/winebarrel/orderedmap)](https://goreportcard.com/report/github.com/winebarrel/orderedmap)

A generic, thread-safe ordered map for Go that preserves insertion order using [container/list](https://pkg.go.dev/container/list).

## Installation

```
go get github.com/winebarrel/orderedmap
```

Requires Go 1.23+.

## Usage

```go
package main

import (
	"fmt"

	"github.com/winebarrel/orderedmap"
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
	pairs := om.Pairs()
	fmt.Println(pairs)
	//=> [{zoo 100} {baz true} {foo new_bar}]

	// Delete a key
	om.Delete("zoo")

	// Get the number of entries
	fmt.Println(om.Len())
	//=> 2

	// Clear all entries
	om.Clear()
}
```
