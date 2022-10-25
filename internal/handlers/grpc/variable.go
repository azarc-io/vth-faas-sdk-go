package grpc

import (
	"context"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
)

type VariableHandler struct {
	client sdk_v1.ManagerServiceClient
}

func NewVariableHandler() sdk_v1.VariableHandler {
	return VariableHandler{}
}

func (g VariableHandler) Set(jobKey string, variables ...*sdk_v1.Variable) error {
	_, err := g.client.SetVariables(context.Background(), sdk_v1.NewSetVariablesRequest(jobKey, variables...))
	return err
}

func (g VariableHandler) Get(jobKey string, names ...string) *sdk_v1.Inputs {
	variables, err := g.client.GetVariables(context.Background(), sdk_v1.NewGetVariablesRequest(jobKey, names...))
	if err != nil {
		return sdk_v1.NewInputs(err)
	}
	var vars []*sdk_v1.Variable
	for _, v := range variables.Variables {
		vars = append(vars, v)
	}
	return sdk_v1.NewInputs(err, vars...)
}
