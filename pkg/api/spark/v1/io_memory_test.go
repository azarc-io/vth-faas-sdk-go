package sdk_v1

import (
	"fmt"
	"testing"
)

type inMemoryIOHandler struct {
	variables map[string]*Variable
	t         *testing.T
}

func NewInMemoryIOHandler(t *testing.T) IOHandler {
	i := &inMemoryIOHandler{t: t, variables: map[string]*Variable{}}
	return i
}

func (i *inMemoryIOHandler) Inputs(jobKey string, names ...string) *Inputs {
	var (
		vars []*Variable
		err  error
	)
	for _, n := range names {
		key := i.key(jobKey, n)
		if v, ok := i.variables[key]; ok {
			var va *Variable
			va, err = NewVariable(v.Name, v.MimeType, v.Value)
			vars = append(vars, va)
		}
	}
	if len(vars) == 0 {
		i.t.Fatalf("no variables found for the params: ")
	}
	return NewInputs(err, vars...)
}

func (i *inMemoryIOHandler) Input(jobKey, name string) *Input {
	inputs := i.Inputs(jobKey, name)
	return inputs.Get(name)
}

func (i *inMemoryIOHandler) Output(jobKey string, variables ...*Variable) error {
	for _, v := range variables {
		i.variables[i.key(jobKey, v.Name)] = v
	}
	return nil
}

func (i *inMemoryIOHandler) key(jobKey, name string) string {
	return fmt.Sprintf("%s_%s", jobKey, name)
}
