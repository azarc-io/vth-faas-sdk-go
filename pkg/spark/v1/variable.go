package sparkv1

import (
	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
)

type Var struct {
	Name     string
	MimeType codec.MimeType
	Value    []byte
}

func NewVar(name string, mimeType codec.MimeType, value any) *Var {
	val, _ := codec.Encode(value)
	return &Var{name, mimeType, val}
}
