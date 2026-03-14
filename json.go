package orderedmap

import (
	"bytes"
	"container/list"
	"encoding"
	"encoding/json"
	"fmt"
)

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
