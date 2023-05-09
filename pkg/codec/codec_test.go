package codec

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{name: "bytes", input: []byte(`my test`), expected: []byte("my test")},
		{name: "raw string", input: "my test string", expected: "my test string"},
		{name: "json int", input: 123, expected: 123.0},
		{name: "json float", input: 123.12, expected: 123.12},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encIn, err := Encode(test.input)
			assert.NoError(t, err)

			var out any
			if _, ok := test.expected.([]byte); ok {
				var o []byte
				err = Decode(encIn, &o)
				out = o
			} else {
				err = Decode(encIn, &out)
			}

			assert.NoError(t, err)
			assert.Equal(t, out, test.expected)
		})
	}
}

func TestIsJson(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		isJson bool
	}{
		{"Open {", " {", true},
		{"Open [", " [", true},
		{"Integer", " 212", true},
		{"Float", " 212.033", true},
		{"True", " true", true},
		{"False", " false", true},
		{"Letters", " hello", false},
		{"Json string", " \"hello", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.isJson, isValidJSON([]byte(test.input)), fmt.Sprintf("Should be Json: %s", test.input))
		})
	}
}
