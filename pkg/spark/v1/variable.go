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
	return &Var{name, mimeType, value, true}
}

type rawVar struct {
	Raw []byte
}
