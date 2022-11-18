package spark_v1

import (
	"fmt"
	"testing"
)

type inMemoryIOHandler struct {
	variables map[string]*Var
	t         *testing.T
}

func NewInMemoryIOHandler(t *testing.T) TestIOHandler {
	i := &inMemoryIOHandler{t: t, variables: map[string]*Var{}}
	return i
}

func (i *inMemoryIOHandler) Inputs(jobKey string, names ...string) Inputs {
	var (
		vars []*Variable
		err  error
	)
	for _, n := range names {
		key := i.key(jobKey, n)
		if v, ok := i.variables[key]; ok {
			var va *Variable
			va, err = newVariable(v.Name, v.MimeType, v.Value)
			vars = append(vars, va)
		}
	}
	if len(vars) == 0 {
		i.t.Fatal("no variables found for the params: ")
	}
	return newInputs(err, vars...)
}

func (i *inMemoryIOHandler) Input(jobKey, name string) Input {
	inputs := i.Inputs(jobKey, name)
	return inputs.Get(name)
}

func (i *inMemoryIOHandler) Output(jobKey string, variables ...*Var) error {
	for _, v := range variables {
		i.variables[i.key(jobKey, v.Name)] = v
	}
	return nil
}

func (i *inMemoryIOHandler) SetVar(jobKey string, v *Var) {
	i.variables[i.key(jobKey, v.Name)] = v
}

func (i *inMemoryIOHandler) key(jobKey, name string) string {
	return fmt.Sprintf("%s_%s", jobKey, name)
}
