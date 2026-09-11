package orderedmap

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"reflect"

	"github.com/winebarrel/linkedlist"
)

// MarshalJSONTo implements [jsonv2.MarshalerTo].
// It encodes the map as a JSON object, preserving insertion order.
func (om *Map[K, V]) MarshalJSONTo(enc *jsontext.Encoder) error {
	om.mu.RLock()
	defer om.mu.RUnlock()

	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}

	if om.entries != nil {
		for e := om.entries.Front(); e != nil; e = e.Next() {
			p := e.Value

			// Keys are written in the object name position,
			// so non-string keys (e.g. int) are stringified by the encoder.
			if err := jsonv2.MarshalEncode(enc, &p.Key); err != nil {
				return err
			}

			if err := jsonv2.MarshalEncode(enc, &p.Value); err != nil {
				return err
			}
		}
	}

	return enc.WriteToken(jsontext.EndObject)
}

// UnmarshalJSONFrom implements [jsonv2.UnmarshalerFrom].
// It decodes a JSON object, preserving the order of its members.
// Entries already held by the map are discarded.
func (om *Map[K, V]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if k := dec.PeekKind(); k != '{' {
		// GoType must be set: the encoding/json v1 wrapper turns this into an
		// UnmarshalTypeError whose Error() dereferences it.
		return &jsonv2.SemanticError{JSONKind: k, GoType: reflect.TypeFor[Map[K, V]]()}
	}

	if _, err := dec.ReadToken(); err != nil {
		return err
	}

	om.mu.Lock()
	defer om.mu.Unlock()

	if om.entries == nil {
		om.entries = linkedlist.New[*Pair[K, V]]()
		om.elementByKey = map[K]*linkedlist.Element[*Pair[K, V]]{}
	} else {
		clear(om.elementByKey)
		om.entries.Init()
	}

	for dec.PeekKind() != '}' {
		var k K
		if err := jsonv2.UnmarshalDecode(dec, &k); err != nil {
			return err
		}

		var v V
		if err := jsonv2.UnmarshalDecode(dec, &v); err != nil {
			return err
		}

		om.set0(k, v)
	}

	_, err := dec.ReadToken()

	return err
}
