package spark_v1

type Var struct {
	Name     string
	MimeType string
	Value    interface{}
}

func NewVar(name, mimeType string, value interface{}) *Var {
	return &Var{name, mimeType, value}
}
