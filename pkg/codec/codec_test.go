package codec

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		mime     MimeType
		expected any
	}{
		{name: "bytes", input: []byte(`my test`), mime: MimeTypeOctetStream, expected: []byte("my test")},
		{name: "raw string", input: "my test string", mime: MimeTypeText, expected: "my test string"},
		{name: "json int", input: 123, mime: MimeTypeText, expected: 123.0},
		{name: "json float", input: 123.12, mime: MimeTypeText, expected: 123.12},
		{name: "image", input: []byte("image bytes"), mime: MimeTypeImageJpeg, expected: []byte("image bytes")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encIn, err := Encode(test.input)
			assert.NoError(t, err)

			var out any
			err = DecodeAndBind(encIn, test.mime, &out)
			assert.NoError(t, err)
			assert.Equal(t, out, test.expected)
		})
	}
}

func TestDecode(t *testing.T) {
	t.Run("raw json", func(t *testing.T) {
		val := `
{
	"foo": "bar"
}
`
		out, err := DecodeValue([]byte(val), MimeTypeJson)
		assert.NoError(t, err)
		assert.Equal(t, map[string]any{"foo": "bar"}, out)
	})

	t.Run("json encoded string", func(t *testing.T) {
		val := `"hello from here"`
		out, err := DecodeValue([]byte(val), MimeTypeText)
		assert.NoError(t, err)
		assert.Equal(t, "hello from here", out)
	})

	t.Run("number", func(t *testing.T) {
		val := `123`
		out, err := DecodeValue([]byte(val), MimeTypeText)
		assert.NoError(t, err)
		assert.Equal(t, 123.0, out)
	})

	t.Run("float", func(t *testing.T) {
		val := `123.567`
		out, err := DecodeValue([]byte(val), MimeTypeText)
		assert.NoError(t, err)
		assert.Equal(t, 123.567, out)
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

func TestType(t *testing.T) {
	tests := []struct {
		mime     MimeType
		expected string
	}{
		{mime: MimeTypeJson, expected: TypeApplication},
		{mime: MimeTypeText, expected: TypeApplication},
		{mime: MimeTypeOctetStream, expected: TypeApplication},
		{mime: MimeTypeJson, expected: TypeApplication},
		{mime: MimeTypeImageJpeg, expected: TypeImage},
		{mime: MimeTypeImagePng, expected: TypeImage},
		{mime: MimeTypeImageGif, expected: TypeImage},
		{mime: MimeTypeImageSvg, expected: TypeImage},
		{mime: MimeTypeImageSvg.WithType("xml"), expected: TypeImage},
		{mime: MimeTypeText.WithType("json"), expected: TypeApplication},
	}

	for _, tc := range tests {
		t.Run(string(tc.mime), func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.mime.Type())
		})
	}
}
