package test

import (
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	sdk_errors "github.com/azarc-io/vth-faas-sdk-go/pkg/errors"
)

// TODO implement the DEFAULT GRPC handler

type inMemoryVariableHandler struct {
	variables map[string]*sdk_v1.Variable
}

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
