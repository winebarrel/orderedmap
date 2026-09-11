package orderedmap_test

import (
	"bytes"
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/orderedmap/v2"
)

// errWriter is an io.Writer that always fails.
type errWriter struct{}

func (w errWriter) Write([]byte) (int, error) {
	return 0, errors.New("write error")
}

func TestMarshalerToInterface(t *testing.T) {
	assert.Implements(t, (*jsonv2.MarshalerTo)(nil), orderedmap.New[string, any]())
	assert.Implements(t, (*jsonv2.UnmarshalerFrom)(nil), orderedmap.New[string, any]())
}

func TestMarshalJSONTo(t *testing.T) {
	tests := []struct {
		name     string
		init     []pair[any]
		expected string
	}{
		{
			name:     "basic object",
			init:     []pair[any]{{k: "foo", v: "bar"}, {k: "zoo", v: 100}, {k: "baz", v: true}},
			expected: `{"foo":"bar","zoo":100,"baz":true}`,
		},
		{
			name:     "empty object",
			init:     []pair[any]{},
			expected: `{}`,
		},
		{
			name:     "nested object",
			init:     []pair[any]{{k: "a", v: map[string]any{"x": 1}}, {k: "b", v: 2}},
			expected: `{"a":{"x":1},"b":2}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := orderedmap.New[string, any]()
			for _, p := range tt.init {
				om.Set(p.k, p.v)
			}
			b, err := jsonv2.Marshal(om)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, string(b))
		})
	}
}

func TestMarshalJSONToPreservesOrder(t *testing.T) {
	om := orderedmap.New[string, any]()
	om.Set("z", 3)
	om.Set("a", 1)
	om.Set("m", 2)

	b, err := jsonv2.Marshal(om)
	assert.NoError(t, err)
	assert.Equal(t, `{"z":3,"a":1,"m":2}`, string(b))
}

func TestMarshalJSONToZeroValueMap(t *testing.T) {
	// Map created by var declaration (not New) has nil internal fields.
	var om orderedmap.Map[string, any]
	b, err := jsonv2.Marshal(&om)
	assert.NoError(t, err)
	assert.Equal(t, `{}`, string(b))
}

func TestMarshalJSONToRespectsOptions(t *testing.T) {
	om := orderedmap.New[string, any]()
	om.Set("a", "<b>")

	// Unlike encoding/json v1, json/v2 does not escape HTML by default.
	b, err := jsonv2.Marshal(om)
	assert.NoError(t, err)
	assert.Equal(t, `{"a":"<b>"}`, string(b))

	b, err = jsonv2.Marshal(om, jsontext.EscapeForHTML(true))
	assert.NoError(t, err)
	assert.Equal(t, `{"a":"\u003cb\u003e"}`, string(b))

	b, err = jsonv2.Marshal(om, jsontext.WithIndent("  "))
	assert.NoError(t, err)
	assert.Equal(t, "{\n  \"a\": \"<b>\"\n}", string(b))
}

func TestMarshalJSONToNestedInStruct(t *testing.T) {
	om := orderedmap.New[string, any]()
	om.Set("z", 1)
	om.Set("a", 2)

	s := struct {
		Map *orderedmap.Map[string, any] `json:"map"`
		Num int                          `json:"num"`
	}{Map: om, Num: 3}

	b, err := jsonv2.Marshal(s)
	assert.NoError(t, err)
	assert.Equal(t, `{"map":{"z":1,"a":2},"num":3}`, string(b))
}

func TestMarshalJSONToTextMarshalerKey(t *testing.T) {
	om := orderedmap.New[textKey, any]()
	om.Set("foo", "bar")
	om.Set("baz", 1)

	b, err := jsonv2.Marshal(om)
	assert.NoError(t, err)
	assert.Equal(t, `{"foo":"bar","baz":1}`, string(b))
}

func TestMarshalJSONToTextMarshalerKeyError(t *testing.T) {
	om := orderedmap.New[errMarshalTextKey, any]()
	om.Set(errMarshalTextKey{}, "val")

	_, err := jsonv2.Marshal(om)
	assert.Error(t, err)
}

func TestMarshalJSONToIntKey(t *testing.T) {
	om := orderedmap.New[int, any]()
	om.Set(1, "one")
	om.Set(2, "two")

	b, err := jsonv2.Marshal(om)
	assert.NoError(t, err)
	assert.Equal(t, `{"1":"one","2":"two"}`, string(b))
}

func TestMarshalJSONToKeyError(t *testing.T) {
	om := orderedmap.New[errMarshalKey, any]()
	om.Set(errMarshalKey{}, "val")

	_, err := jsonv2.Marshal(om)
	assert.Error(t, err)
}

func TestMarshalJSONToValueError(t *testing.T) {
	om := orderedmap.New[string, errMarshalValue]()
	om.Set("key", errMarshalValue{})

	_, err := jsonv2.Marshal(om)
	assert.Error(t, err)
}

func TestMarshalJSONToEncoderWriteError(t *testing.T) {
	om := orderedmap.New[string, any]()
	om.Set("foo", "bar")

	err := jsonv2.MarshalWrite(errWriter{}, om)
	assert.Error(t, err)
}

func TestMarshalJSONToEncode(t *testing.T) {
	om := orderedmap.New[string, any]()
	om.Set("foo", "bar")

	var buf bytes.Buffer
	enc := jsontext.NewEncoder(&buf)
	assert.NoError(t, jsonv2.MarshalEncode(enc, om))
	assert.NoError(t, jsonv2.MarshalEncode(enc, om))
	assert.Equal(t, "{\"foo\":\"bar\"}\n{\"foo\":\"bar\"}\n", buf.String())
}

func TestUnmarshalJSONFrom(t *testing.T) {
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
		{
			name:    "truncated key",
			input:   `{"ke`,
			wantErr: true,
		},
		{
			name:    "empty input",
			input:   ``,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := orderedmap.New[string, any]()
			err := jsonv2.Unmarshal([]byte(tt.input), om)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, mapToPairs(t, om))
		})
	}
}

func TestUnmarshalJSONFromPreservesOrder(t *testing.T) {
	om := orderedmap.New[string, any]()
	err := jsonv2.Unmarshal([]byte(`{"z":"last","a":"first","m":"middle"}`), om)
	assert.NoError(t, err)
	assert.Equal(t, []pair[any]{{k: "z", v: "last"}, {k: "a", v: "first"}, {k: "m", v: "middle"}}, mapToPairs(t, om))
}

func TestUnmarshalJSONFromOverwritesExisting(t *testing.T) {
	om := orderedmap.New[string, any]()
	om.Set("old", "value")

	err := jsonv2.Unmarshal([]byte(`{"new":"data"}`), om)
	assert.NoError(t, err)
	assert.Equal(t, []pair[any]{{k: "new", v: "data"}}, mapToPairs(t, om))
}

func TestUnmarshalJSONFromNilMap(t *testing.T) {
	// Map created by var declaration (not New) has nil internal fields.
	var om orderedmap.Map[string, any]
	err := jsonv2.Unmarshal([]byte(`{"foo":"bar"}`), &om)
	assert.NoError(t, err)
	assert.Equal(t, []pair[any]{{k: "foo", v: "bar"}}, mapToPairs(t, &om))
}

func TestUnmarshalJSONFromDuplicateNames(t *testing.T) {
	// json/v2 rejects duplicate object names by default.
	om := orderedmap.New[string, any]()
	err := jsonv2.Unmarshal([]byte(`{"a":1,"a":2}`), om)
	assert.Error(t, err)

	om = orderedmap.New[string, any]()
	err = jsonv2.Unmarshal([]byte(`{"a":1,"a":2}`), om, jsontext.AllowDuplicateNames(true))
	assert.NoError(t, err)
	assert.Equal(t, []pair[any]{{k: "a", v: float64(2)}}, mapToPairs(t, om))
}

func TestUnmarshalJSONFromTextUnmarshalerKey(t *testing.T) {
	om := orderedmap.New[textKey, any]()
	err := jsonv2.Unmarshal([]byte(`{"foo":"bar","baz":true}`), om)
	assert.NoError(t, err)
	assert.Equal(t, []orderedmap.Pair[textKey, any]{{Key: "foo", Value: "bar"}, {Key: "baz", Value: true}}, om.Entries())
}

func TestUnmarshalJSONFromTextUnmarshalerKeyError(t *testing.T) {
	om := orderedmap.New[errTextKey, any]()
	err := jsonv2.Unmarshal([]byte(`{"foo":"bar"}`), om)
	assert.Error(t, err)
}

func TestUnmarshalJSONFromIntKey(t *testing.T) {
	om := orderedmap.New[int, any]()
	err := jsonv2.Unmarshal([]byte(`{"1":1}`), om)
	assert.NoError(t, err)
	assert.Equal(t, []orderedmap.Pair[int, any]{{Key: 1, Value: float64(1)}}, om.Entries())
}

func TestUnmarshalJSONFromIntKeyError(t *testing.T) {
	// K=int cannot accept a non-numeric string key from JSON.
	om := orderedmap.New[int, any]()
	err := jsonv2.Unmarshal([]byte(`{"abc":1}`), om)
	assert.Error(t, err)
}

func TestUnmarshalJSONFromValueError(t *testing.T) {
	om := orderedmap.New[string, errValue]()
	err := jsonv2.Unmarshal([]byte(`{"key":{}}`), om)
	assert.Error(t, err)
}

func TestUnmarshalJSONFromDecode(t *testing.T) {
	dec := jsontext.NewDecoder(bytes.NewReader([]byte(`{"z":1,"a":2} {"m":3}`)))

	om := orderedmap.New[string, any]()
	assert.NoError(t, jsonv2.UnmarshalDecode(dec, om))
	assert.Equal(t, []pair[any]{{k: "z", v: float64(1)}, {k: "a", v: float64(2)}}, mapToPairs(t, om))

	assert.NoError(t, jsonv2.UnmarshalDecode(dec, om))
	assert.Equal(t, []pair[any]{{k: "m", v: float64(3)}}, mapToPairs(t, om))
}

func TestJSONv2RoundTrip(t *testing.T) {
	om := orderedmap.New[string, any]()
	om.Set("z", "last")
	om.Set("a", float64(1))
	om.Set("m", true)
	om.Set(`foo\bar"baz`, "tab\there\nnew\x00line")

	b, err := jsonv2.Marshal(om)
	assert.NoError(t, err)

	om2 := orderedmap.New[string, any]()
	err = jsonv2.Unmarshal(b, om2)
	assert.NoError(t, err)
	assert.Equal(t, om.Entries(), om2.Entries())
}

func TestMarshalJSONToObjectNamePosition(t *testing.T) {
	// A JSON object cannot be written where an object name is expected.
	var buf bytes.Buffer
	enc := jsontext.NewEncoder(&buf)
	assert.NoError(t, enc.WriteToken(jsontext.BeginObject))

	om := orderedmap.New[string, any]()
	assert.Error(t, om.MarshalJSONTo(enc))
}

func TestUnmarshalJSONFromObjectNamePosition(t *testing.T) {
	// A JSON object cannot be read where an object name is expected.
	dec := jsontext.NewDecoder(bytes.NewReader([]byte(`{{`)))
	_, err := dec.ReadToken()
	assert.NoError(t, err)

	om := orderedmap.New[string, any]()
	assert.Error(t, om.UnmarshalJSONFrom(dec))
}

func TestUnmarshalJSONFromNonObjectErrorMessage(t *testing.T) {
	for _, input := range []string{`null`, `[1,2]`, `"x"`, `1`} {
		t.Run(input, func(t *testing.T) {
			err := jsonv2.Unmarshal([]byte(input), orderedmap.New[string, int]())
			assert.ErrorContains(t, err, "orderedmap.Map[string,int]")
		})
	}
}
