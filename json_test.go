package orderedmap_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/orderedmap"
)

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
	// JSON object keys should be inserted in the order they appear
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
