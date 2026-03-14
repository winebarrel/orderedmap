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

		var keyStr string
		if tm, ok := any(p.Key).(encoding.TextMarshaler); ok {
			text, err := tm.MarshalText()
			if err != nil {
				return nil, fmt.Errorf("orderedmap: cannot marshal key: %w", err)
			}
			keyStr = string(text)
		} else {
			keyJSON, err := json.Marshal(p.Key)
			if err != nil {
				return nil, fmt.Errorf("orderedmap: cannot marshal key: %w", err)
			}
			// If the key marshaled as a JSON string, use it directly; otherwise wrap it.
			if len(keyJSON) > 0 && keyJSON[0] == '"' {
				keyStr = string(keyJSON[1 : len(keyJSON)-1])
			} else {
				keyStr = string(keyJSON)
			}
		}

		keyBytes, _ := json.Marshal(keyStr)
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
