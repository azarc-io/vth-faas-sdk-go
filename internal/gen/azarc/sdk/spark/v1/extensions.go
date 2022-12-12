package sparkv1

import (
	jsoniter "github.com/json-iterator/go"
)

const (
	NoMimeType   = ""
	MimeTypeJSON = "application/json"
)

var json2 = jsoniter.ConfigFastest

/************************************************************************/
// MARSHALLING
/************************************************************************/

type serdes struct {
	Unmarshal func(value []byte, a any) error
	Marshal   func(a any) ([]byte, error)
}

var SerdesMap = map[string]serdes{
	MimeTypeJSON: {
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
	NoMimeType: {
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
	return SerdesMap[MimeTypeJSON].Unmarshal(x.Data, a)
}
