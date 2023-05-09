package codec

import (
	"encoding/base64"
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

func TestDecodeToBytes(t *testing.T) {
	t.Run("raw json", func(t *testing.T) {
		val := `
{
	"foo": "bar
}
`
		enc := base64.StdEncoding.EncodeToString([]byte(val))
		out, err := DecodeToBytes(enc)
		assert.NoError(t, err)
		assert.Equal(t, val, string(out))
	})

	t.Run("raw string", func(t *testing.T) {
		val := `hello from here`
		enc := base64.StdEncoding.EncodeToString([]byte(val))
		out, err := DecodeToBytes(enc)
		assert.NoError(t, err)
		assert.Equal(t, val, string(out))
	})

	t.Run("byte array", func(t *testing.T) {
		val := `hello from here`
		enc := []byte(val)
		out, err := DecodeToBytes(enc)
		assert.NoError(t, err)
		assert.Equal(t, val, string(out))
	})
}

func TestMime(t *testing.T) {
	t.Run("Get Base Type", func(t *testing.T) {
		mt := MimeTypeOctetStream.WithType("pdf")
		assert.Equal(t, MimeTypeOctetStream, mt.BaseType())
	})

	t.Run("With sub type", func(t *testing.T) {
		mt := MimeTypeOctetStream.WithType("pdf")
		assert.Equal(t, MimeType("application/octet-stream+pdf"), mt)
	})
}
