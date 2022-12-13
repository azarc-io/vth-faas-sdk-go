package spark_v1

type Var struct {
	Name     string
	MimeType string
	Value    any
}

func NewVar(name, mimeType string, value any) *Var {
	return &Var{name, mimeType, value}
}
