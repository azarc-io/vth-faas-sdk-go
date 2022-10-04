package sdk_v1

import (
	"encoding/json"
)

// TODO yaml, xml, json, toml, csv <- just the ones that golang support
// we need to use the correct encoder based on the mime_type field of the message

func (x *Variable) Raw() []byte {
	return x.GetValue()
}

func (x *Variable) Bind(a any) error {
	return json.Unmarshal(x.Value, a)
}

func (x *Stage) Raw() []byte {
	return nil
}

func (x *Stage) Bind(a any) error {
	return nil
}

func NewStageResult(stage *Stage, data any) *StageResult {
	b, err := json.Marshal(data)
	if err != nil {
		panic("this should panic?") // TODO proper error handling
	}
	return &StageResult{
		Stage: stage,
		Data:  b,
	}
}

func NewVariable(name, mimeType string, value []byte) *Variable {
	return &Variable{
		Name:     name,
		Value:    value,
		MimeType: mimeType,
	}
}
