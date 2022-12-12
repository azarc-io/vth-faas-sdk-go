package spark_v1

import (
	"encoding/json"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/internal/gen/azarc/sdk/spark/v1"
)

func MarshalBinary(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func UnmarshalBinaryTo(data []byte, out interface{}, mimeType string) error {
	if mimeType == "" {
		return sparkv1.SerdesMap[MimeTypeJSON].Unmarshal(data, &out)
	} else {
		return sparkv1.SerdesMap[mimeType].Unmarshal(data, &out)
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
