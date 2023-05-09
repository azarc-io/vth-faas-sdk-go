package codec

import (
	"encoding/base64"
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

	return json.Unmarshal(input, target)
}

func DecodeToBytes(input any) ([]byte, error) {
	if input == nil {
		return nil, nil
	}

	var err error
	switch val := input.(type) {
	case string:
		input, err = base64.StdEncoding.DecodeString(val)
		if err != nil {
			return nil, err
		}
	}

	data, ok := input.([]byte)
	if !ok {
		return nil, ErrInvalidOctetStreamType
	}

	// check if this is a JSON string instead
	var nv string
	if err2 := json.Unmarshal(data, &nv); err2 == nil {
		// covers test: "json bytes from raw json string"
		data = []byte(nv)
	}

	// Check if this is double encoded byte array
	if IsBase64(string(data)) {
		if dbl, err := base64.StdEncoding.DecodeString(string(data)); err == nil {
			data = dbl
		}
	}

	return data, nil
}

// IsBase64 checks if a string is a valid base64 encoded string
func IsBase64(s string) bool {
	if len(s)%4 != 0 || !validBase64.MatchString(s) {
		return false
	}

	return true
}
