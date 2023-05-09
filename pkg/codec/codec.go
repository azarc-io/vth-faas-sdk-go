package codec

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var (
	ErrTargetNotPointer        = errors.New("unable to set value of non-pointer")
	ErrInvalidOctetStreamType  = errors.New("unable to unmarshal octet-stream type")
	ErrInvalidJsonType         = errors.New("unable to unmarshal json type")
	ErrUnableMarshalMimeType   = errors.New("unable to marshal with unknown mime type")
	ErrUnableUnmarshalMimeType = errors.New("unable to unmarshal with unknown mime type")
)

type MimeType string

const (
	MimeTypeJson        MimeType = "application/json"
	MimeTypeText        MimeType = "application/text"
	MimeTypeOctetStream MimeType = "application/octet-stream"
)

func (mt MimeType) WithType(subType string) MimeType {
	return MimeType(fmt.Sprintf("%s+%s", mt, strings.ToLower(subType)))
}

func Encode(v any) ([]byte, error) {
	if val, ok := v.([]byte); ok {
		return val, nil
	}

	return json.Marshal(v)
}

func Decode(input []byte, target any) error {
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return ErrTargetNotPointer
	}

	if len(input) == 0 {
		return nil
	}

	// this is fallback for non-json types being set
	switch target.(type) {
	case *[]byte:
		rv.Elem().Set(reflect.ValueOf(input))
		return nil
	}

	return json.Unmarshal(input, target)
}

var jsonRegex = regexp.MustCompile(`^\s*(\{|\[|"|[0-9.]+|true|false)`)

func isValidJSON(data []byte) bool {
	v := jsonRegex.FindSubmatchIndex(data)
	return len(v) > 0
}
