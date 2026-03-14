package orderedmap

import (
	"bytes"
	"container/list"
	"encoding"
	"encoding/json"
	"fmt"
)

func (om *Map[K, V]) MarshalJSON() ([]byte, error) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	var buf bytes.Buffer
	buf.WriteByte('{')

	first := true
	for e := om.entries.Front(); e != nil; e = e.Next() {
		if !first {
			buf.WriteByte(',')
		}
		first = false

		p := e.Value.(*Pair[K, V])

		var keyBytes []byte
		if tm, ok := any(p.Key).(encoding.TextMarshaler); ok {
			text, err := tm.MarshalText()
			if err != nil {
				return nil, fmt.Errorf("orderedmap: cannot marshal key: %w", err)
			}
			keyBytes, _ = json.Marshal(string(text))
		} else {
			var err error
			keyBytes, err = json.Marshal(p.Key)
			if err != nil {
				return nil, fmt.Errorf("orderedmap: cannot marshal key: %w", err)
			}
			// Non-string keys (e.g. int) must be wrapped as JSON strings for object keys.
			if len(keyBytes) > 0 && keyBytes[0] != '"' {
				keyBytes, _ = json.Marshal(string(keyBytes))
			}
		}
		buf.Write(keyBytes)
		buf.WriteByte(':')

		valBytes, err := json.Marshal(p.Value)
		if err != nil {
			return nil, fmt.Errorf("orderedmap: cannot marshal value: %w", err)
		}
		buf.Write(valBytes)
	}

	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func (om *Map[K, V]) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))

	t, err := dec.Token()
	if err != nil {
		return err
	}

	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return fmt.Errorf("orderedmap: expected JSON object, got %v", t)
	}

	om.mu.Lock()
	defer om.mu.Unlock()

	if om.entries == nil {
		om.entries = list.New()
		om.elementByKey = map[K]*list.Element{}
	} else {
		clear(om.elementByKey)
		om.entries.Init()
	}

	for dec.More() {
		kt, err := dec.Token()
		if err != nil {
			return err
		}

		keyStr := kt.(string) // JSON object keys are always strings

		var k K
		if tu, ok := any(&k).(encoding.TextUnmarshaler); ok {
			if err := tu.UnmarshalText([]byte(keyStr)); err != nil {
				return fmt.Errorf("orderedmap: cannot unmarshal key %q: %w", keyStr, err)
			}
		} else {
			keyJSON, _ := json.Marshal(keyStr)
			if err := json.Unmarshal(keyJSON, &k); err != nil {
				// Fallback: try unmarshaling the key as a raw JSON value (e.g. numeric keys).
				if err2 := json.Unmarshal([]byte(keyStr), &k); err2 != nil {
					return fmt.Errorf("orderedmap: cannot unmarshal key %q into %T: %w", keyStr, k, err)
				}
			}
		}

		var v V
		if err := dec.Decode(&v); err != nil {
			return err
		}

		om.set0(k, v)
	}

	_, err = dec.Token() // closing }
	return err
}
