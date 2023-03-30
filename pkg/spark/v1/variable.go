package sparkv1

import (
	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
)

type Var struct {
	Name     string
	MimeType codec.MimeType
	Value    any
	Raw      bool
}

func NewVar(name string, mimeType codec.MimeType, value any) *Var {
	return &Var{name, mimeType, value, false}
}

func NewRawVar(name string, mimeType codec.MimeType, value []byte) *Var {
	//resp := base64.StdEncoding.EncodeToString(value)
	resp := value
	return &Var{name, mimeType, resp, true}
}

type rawVar struct {
	Raw []byte
}
