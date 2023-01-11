package sparkv1

import "github.com/azarc-io/vth-faas-sdk-go/pkg/codec"

type Var struct {
	Name     string
	MimeType codec.MimeType
	Value    any
}

func NewVar(name string, mimeType codec.MimeType, value any) *Var {
	return &Var{name, mimeType, value}
}
