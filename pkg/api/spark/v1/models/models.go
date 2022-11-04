package models

type Variable struct {
	Name     string
	MimeType string
	Value    any
}

func NewVar(name, mimeType string, value any) *Variable {
	return &Variable{name, mimeType, value}
}
