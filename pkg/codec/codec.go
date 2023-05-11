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

// isBase64Characters checks if a string contains only valid base64 characters
var validBase64 = regexp.MustCompile(`^[A-Za-z0-9+/]*=?=?$`)

func (mt MimeType) WithType(subType string) MimeType {
	return MimeType(fmt.Sprintf("%s+%s", mt, strings.ToLower(subType)))
}

func (mt MimeType) BaseType() MimeType {
	v := strings.Split(string(mt), "+")[0]
	return MimeType(v)
}

func Encode(v any) ([]byte, error) {
	if b, ok := v.([]byte); ok {
		return b, nil
	}
	return json.Marshal(v)
}

func DecodeValue(input []byte, mime MimeType) (any, error) {
	var data any
	return data, DecodeAndBind(input, mime, &data)
}

func DecodeAndBind(input []byte, mime MimeType, target any) error {
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Pointer || val.IsNil() {
		return fmt.Errorf("DecodeAndBind: expected a pointer")
	}

	elem := val.Elem()
	if !elem.CanSet() {
		return fmt.Errorf("overrideValue: cannot set value of the pointer")
	}

	// allow binding of raw bytes
	_, ok := target.(*[]byte)
	if mime.BaseType() == MimeTypeOctetStream || ok {
		elem.Set(reflect.ValueOf(input))
		return nil
	}

	return json.Unmarshal(input, target)
}
