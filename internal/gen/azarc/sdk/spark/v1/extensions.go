package sparkv1

import (
	"encoding/json"
	"github.com/azarc-io/vth-faas-sdk-go/internal/common"
	jsoniter "github.com/json-iterator/go"
	"reflect"
)

// TODO get rid of this whole extension in the future

var json2 = jsoniter.ConfigFastest

/************************************************************************/
// MARSHALLING
/************************************************************************/

type serdes struct {
	Unmarshal func(value []byte, a any) error
	Marshal   func(a any) ([]byte, error)
}

var SerdesMap = map[string]serdes{
	common.MimeTypeJSON: {
		Unmarshal: func(value []byte, a any) error {
			return json2.Unmarshal(value, a)
		},
		Marshal: func(a any) ([]byte, error) {
			b, err := json2.Marshal(a)
			if err != nil {
				return nil, err
			}
			return b, nil
		},
	},
	common.NoMimeType: {
		Unmarshal: func(value []byte, a any) error {
			return json2.Unmarshal(value, a)
		},
		Marshal: func(a any) ([]byte, error) {
			b, err := json2.Marshal(a)
			if err != nil {
				return nil, err
			}
			return b, nil
		},
	},
}

/************************************************************************/
// VARIABLE EXTENSIONS
/************************************************************************/

func (x *Variable) Raw() ([]byte, error) {
	return x.Data, nil
}

func (x *Variable) Bind(a any) error {
	return SerdesMap[x.MimeType].Unmarshal(x.Data, a)
}

/************************************************************************/
// STAGE RESULT EXTENSIONS
/************************************************************************/

func (x *GetStageResultResponse) Raw() ([]byte, error) {
	return x.Data, nil
}

func (x *GetStageResultResponse) Bind(a any) error {
	return SerdesMap[common.MimeTypeJSON].Unmarshal(x.Data, a)
}

/************************************************************************/
// HELPERS
/************************************************************************/

func MarshalBinary(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func UnmarshalBinaryTo(data []byte, out interface{}, mimeType string) error {
	switch mimeType {
	case common.MimeTypeJSON:
		return SerdesMap[mimeType].Unmarshal(data, &out)
	default:
		v := reflect.ValueOf(out).Elem()
		v.SetString(string(data))
		return nil
	}
}

func ConvertBytes(data []byte, mimeType string) (out []byte, err error) {
	var value interface{}
	err = UnmarshalBinaryTo(data, &value, mimeType)
	if err != nil {
		return
	}

	switch v := value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return data, nil
	case string:
		return []byte(v), nil
	default:
		err = UnmarshalBinaryTo(data, &out, mimeType)
	}

	return
}
