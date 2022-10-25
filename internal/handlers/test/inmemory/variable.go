package inmemory

import (
	"fmt"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	"testing"
)

type inMemoryVariableHandler struct {
	variables map[string]*sdk_v1.Variable
	t         *testing.T
}

func NewVariableHandler(t *testing.T, variables *sdk_v1.SetVariablesRequest) sdk_v1.VariableHandler {
	i := &inMemoryVariableHandler{t: t, variables: make(map[string]*sdk_v1.Variable)}
	for k, v := range variables.GetVariables() {
		i.variables[i.key(k, variables.JobKey)] = v
	}
	return i
}

func (i *inMemoryVariableHandler) Get(jobKey string, names ...string) *sdk_v1.Inputs {
	var vars []*sdk_v1.Variable
	for _, n := range names {
		key := i.key(n, jobKey)
		if v, ok := i.variables[key]; ok {
			vars = append(vars, v)
		}
	}
	if len(vars) == 0 {
		i.t.Fatalf("no variables found for the params: ")
	}
	return sdk_v1.NewInputs(nil, vars...)
}

func (i *inMemoryVariableHandler) Set(jobKey string, variables ...*sdk_v1.Variable) error {
	for _, v := range variables {
		i.variables[i.key(v.Name, jobKey)] = v
	}
	return nil
}

func (i *inMemoryVariableHandler) key(name, jobKey string) string {
	return fmt.Sprintf("%s_%s", name, jobKey)
}
