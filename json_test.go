package orderedmap_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/orderedmap"
)

// textKey is a string key type that implements encoding.TextUnmarshaler.
type textKey string

func (k *textKey) UnmarshalText(b []byte) error {
	*k = textKey(b)
	return nil
}

// errTextKey is a key type whose TextUnmarshaler always returns an error.
type errTextKey struct{}

func (k *errTextKey) UnmarshalText([]byte) error {
	return errors.New("text unmarshal error")
}

// errValue is a value type whose JSON unmarshaling always returns an error.
type errValue struct{}

func (v *errValue) UnmarshalJSON([]byte) error {
	return errors.New("value unmarshal error")
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []pair[any]
		wantErr  bool
	}{
		{
			name:     "basic object",
			input:    `{"foo":"bar","zoo":100,"baz":true}`,
			expected: []pair[any]{{k: "foo", v: "bar"}, {k: "zoo", v: float64(100)}, {k: "baz", v: true}},
		},
		{
			name:     "empty object",
			input:    `{}`,
			expected: []pair[any]{},
		},
		{
			name:     "nested object",
			input:    `{"a":1,"b":{"x":2,"y":3},"c":4}`,
			expected: []pair[any]{{k: "a", v: float64(1)}, {k: "b", v: map[string]any{"x": float64(2), "y": float64(3)}}, {k: "c", v: float64(4)}},
		},
		{
			name:    "non-object",
			input:   `[1,2,3]`,
			wantErr: true,
		},
		{
			name:    "invalid json",
			input:   `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := orderedmap.New[string, any]()
			err := json.Unmarshal([]byte(tt.input), om)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, mapToPairs(t, om))
		})
	}
}

func TestUnmarshalJSONPreservesOrder(t *testing.T) {
	input := `{"z":"last","a":"first","m":"middle"}`
	om := orderedmap.New[string, any]()
	err := json.Unmarshal([]byte(input), om)
	assert.NoError(t, err)
	assert.Equal(t, []pair[any]{{k: "z", v: "last"}, {k: "a", v: "first"}, {k: "m", v: "middle"}}, mapToPairs(t, om))
}

func TestUnmarshalJSONOverwritesExisting(t *testing.T) {
	om := orderedmap.New[string, any]()
	om.Set("old", "value")

	err := json.Unmarshal([]byte(`{"new":"data"}`), om)
	assert.NoError(t, err)
	assert.Equal(t, []pair[any]{{k: "new", v: "data"}}, mapToPairs(t, om))
}

func TestUnmarshalJSONNilMap(t *testing.T) {
	// Map created by var declaration (not New) has nil internal fields.
	var om orderedmap.Map[string, any]
	err := json.Unmarshal([]byte(`{"foo":"bar"}`), &om)
	assert.NoError(t, err)
	assert.Equal(t, []pair[any]{{k: "foo", v: "bar"}}, mapToPairs(t, &om))
}

func TestUnmarshalJSONInitialTokenError(t *testing.T) {
	// Call UnmarshalJSON directly with empty input to trigger dec.Token() error.
	om := orderedmap.New[string, any]()
	err := om.UnmarshalJSON([]byte{})
	assert.Error(t, err)
}

func TestUnmarshalJSONInnerTokenError(t *testing.T) {
	// Call UnmarshalJSON directly with a truncated key to trigger inner dec.Token() error.
	om := orderedmap.New[string, any]()
	err := om.UnmarshalJSON([]byte(`{"ke`))
	assert.Error(t, err)
}

func TestUnmarshalJSONTextUnmarshalerKey(t *testing.T) {
	om := orderedmap.New[textKey, any]()
	err := json.Unmarshal([]byte(`{"foo":"bar","baz":true}`), om)
	assert.NoError(t, err)
	entries := om.Entries()
	assert.Len(t, entries, 2)
	assert.Equal(t, textKey("foo"), entries[0].Key)
	assert.Equal(t, "bar", entries[0].Value)
	assert.Equal(t, textKey("baz"), entries[1].Key)
	assert.Equal(t, true, entries[1].Value)
}

func TestUnmarshalJSONTextUnmarshalerKeyError(t *testing.T) {
	om := orderedmap.New[errTextKey, any]()
	err := json.Unmarshal([]byte(`{"foo":"bar"}`), om)
	assert.Error(t, err)
}

func TestUnmarshalJSONUnmarshalKeyError(t *testing.T) {
	// K=int cannot accept a non-numeric string key from JSON.
	om := orderedmap.New[int, any]()
	err := om.UnmarshalJSON([]byte(`{"abc":1}`))
	assert.Error(t, err)
}

func TestUnmarshalJSONDecodeValueError(t *testing.T) {
	om := orderedmap.New[string, errValue]()
	err := json.Unmarshal([]byte(`{"key":{}}`), om)
	assert.Error(t, err)
}
