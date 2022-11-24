package sparkv1

import (
	"encoding/json"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	NoMimeType   = ""
	MimeTypeJSON = "application/json"
)

/************************************************************************/
// MARSHALLING
/************************************************************************/

type serdes struct {
	Unmarshal func(value *structpb.Value, a any) error
	Marshal   func(a any) (*structpb.Value, error)
}

var SerdesMap = map[string]serdes{
	MimeTypeJSON: {
		Unmarshal: func(value *structpb.Value, a any) error {
			data, err := value.MarshalJSON()
			if err != nil {
				return err
			}
			return json.Unmarshal(data, a)
		},
		Marshal: func(a any) (*structpb.Value, error) {
			value, err := structpb.NewValue(a)
			if err != nil {
				b, err := json.Marshal(a)
				if err != nil {
					return nil, err
				}
				v := map[string]any{}
				err = json.Unmarshal(b, &v)
				if err != nil {
					return nil, err
				}
				return structpb.NewValue(v)
			}
			return value, nil
		},
	},
	NoMimeType: {
		Unmarshal: func(value *structpb.Value, a any) error {
			data, err := value.MarshalJSON()
			if err != nil {
				return err
			}
			return json.Unmarshal(data, a)
		},
		Marshal: func(a any) (*structpb.Value, error) {
			value, err := structpb.NewValue(a)
			if err != nil {
				b, err := json.Marshal(a)
				if err != nil {
					return nil, err
				}
				v := map[string]any{}
				err = json.Unmarshal(b, &v)
				if err != nil {
					return nil, err
				}
				return structpb.NewValue(v)
			}
			return value, nil
		},
	},
}

func GetRawFromPb(data *structpb.Value) ([]byte, error) {
	switch data.Kind.(type) {
	case *structpb.Value_NullValue:
		return nil, nil
	case *structpb.Value_StringValue:
		return []byte(data.GetStringValue()), nil
	}

	return data.MarshalJSON()
}

/************************************************************************/
// VARIABLE EXTENSIONS
/************************************************************************/

func (x *Variable) Raw() ([]byte, error) {
	return x.Value.MarshalJSON()
}

func (x *Variable) Bind(a any) error {
	return SerdesMap[x.MimeType].Unmarshal(x.Value, a)
}

/************************************************************************/
// STAGE RESULT EXTENSIONS
/************************************************************************/

func (x *StageResult) Raw() ([]byte, error) {
	return x.Data.MarshalJSON()
}

func (x *StageResult) Bind(a any) error {
	return SerdesMap[MimeTypeJSON].Unmarshal(x.Data, a)
}
