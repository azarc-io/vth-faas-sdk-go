package test

import (
	sdk_errors "github.com/azarc-io/vth-faas-sdk-go/errors"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
)

type inMemoryVariableHandler struct {
	variables map[string]*sdk_v1.Variable
}

// TODO add t *testing.T
func NewMockVariableHandler() api.VariableHandler {
	return inMemoryVariableHandler{variables: make(map[string]*sdk_v1.Variable)}
}

func (i inMemoryVariableHandler) Get(name string) (*sdk_v1.Variable, error) {
	if variable, ok := i.variables[name]; ok {
		return variable, nil
	}
	return nil, sdk_errors.VariableNotFound
}

func (i inMemoryVariableHandler) Set(variable *sdk_v1.Variable) error {
	i.variables[variable.Name] = variable
	return nil
}
