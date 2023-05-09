package codec_test

import (
	"encoding/base64"
	"fmt"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		ot       reflect.Type
		expected any
	}{
		{name: "raw string to bytes", input: "my test", ot: reflect.TypeOf([]byte{}), expected: []byte("my test")},
		{name: "json bytes", input: fmt.Sprintf(`"%s"`, base64.StdEncoding.EncodeToString([]byte("my test"))), ot: reflect.TypeOf([]byte{}), expected: []byte("my test")},
		{name: "raw string", input: "my test string", ot: reflect.TypeOf(""), expected: "my test string"},
		{name: "json string", input: `"my test string"`, ot: reflect.TypeOf(""), expected: "my test string"},
		{name: "json int", input: `123`, ot: reflect.TypeOf(0), expected: 123.0},
		{name: "json float", input: `123.12`, ot: reflect.TypeOf(0), expected: 123.12},
		{name: "json bytes from raw json string", input: `"{\n  \t\"cntr_no\": \"MSMU6298516\",\n  \t\"carrier_no\": \"MSCU\"\n}"`, ot: reflect.TypeOf([]byte{}), expected: []byte("{\n  \t\"cntr_no\": \"MSMU6298516\",\n  \t\"carrier_no\": \"MSCU\"\n}")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := base64.StdEncoding.EncodeToString([]byte(test.input))

			switch test.ot {
			case reflect.TypeOf(""):
				var out string
				assert.NoError(t, codec.Decode(input, &out))
				assert.Equal(t, test.expected, out)
			case reflect.TypeOf([]byte{}):
				var out []byte
				assert.NoError(t, codec.Decode(input, &out))
				assert.Equal(t, string(test.expected.([]byte)), string(out))
			case reflect.TypeOf(0):
				var out float64
				assert.NoError(t, codec.Decode(input, &out))
				assert.Equal(t, test.expected, out)
			default:
				t.Fatalf("type not supported")
			}
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
		out, err := codec.DecodeToBytes(enc)
		assert.NoError(t, err)
		assert.Equal(t, val, string(out))
	})

	t.Run("raw string", func(t *testing.T) {
		val := `hello from here`
		enc := base64.StdEncoding.EncodeToString([]byte(val))
		out, err := codec.DecodeToBytes(enc)
		assert.NoError(t, err)
		assert.Equal(t, val, string(out))
	})

	t.Run("byte array", func(t *testing.T) {
		val := `hello from here`
		enc := []byte(val)
		out, err := codec.DecodeToBytes(enc)
		assert.NoError(t, err)
		assert.Equal(t, val, string(out))
	})
}
