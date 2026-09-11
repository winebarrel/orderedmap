package orderedmap

import (
	jsonv1 "encoding/json"
	jsonv2 "encoding/json/v2"
)

// MarshalJSON implements [jsonv1.Marshaler] with encoding/json v1 semantics.
// Encoders aware of [jsonv2.MarshalerTo] call [Map.MarshalJSONTo] instead.
func (om *Map[K, V]) MarshalJSON() ([]byte, error) {
	return jsonv2.Marshal(om, jsonv1.DefaultOptionsV1())
}

// UnmarshalJSON implements [jsonv1.Unmarshaler] with encoding/json v1 semantics.
// Decoders aware of [jsonv2.UnmarshalerFrom] call [Map.UnmarshalJSONFrom] instead.
func (om *Map[K, V]) UnmarshalJSON(data []byte) error {
	return jsonv2.Unmarshal(data, om, jsonv1.DefaultOptionsV1())
}
