package codec

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
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

func Encode(v any, mime MimeType) ([]byte, error) {
	if mime == MimeTypeJson.WithType("text") {
		switch val := v.(type) {
		case string:
			return []byte(val), nil
		case []byte:
			return val, nil
		}
	}

	return json.Marshal(v)
}

func Decode(input any, target any) error {
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return ErrTargetNotPointer
	}

	if input == nil {
		return nil
	}

	var err error
	switch val := input.(type) {
	case string:
		input, err = base64.StdEncoding.DecodeString(val)
		if err != nil {
			return err
		}
	}

	data, ok := input.([]byte)
	if !ok {
		return ErrInvalidOctetStreamType
	}

	if err := json.Unmarshal(data, target); err != nil {
		// this is fallback for non-json types being set
		switch target.(type) {
		case *string:
			// is a string
			rv.Elem().Set(reflect.ValueOf(string(data)))
		case *[]byte:
			rv.Elem().Set(reflect.ValueOf(data))
		default:
			return err
		}
	}

	return nil
}
