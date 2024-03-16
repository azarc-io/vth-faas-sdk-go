package module_test_runner

import (
	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
)

func MustEncode(data any) []byte {
	d, err := codec.Encode(data)
	if err != nil {
		panic(err)
	}

	return d
}
