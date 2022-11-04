package grpc

import (
	"context"

	"github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1/models"

	v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1"
)

type VariableHandler struct {
	client v1.ManagerServiceClient
}

func NewIOHandler(client v1.ManagerServiceClient) v1.IOHandler {
	return VariableHandler{client}
}

func (g VariableHandler) Inputs(jobKey string, names ...string) *v1.Inputs {
	variables, err := g.client.GetVariables(context.Background(), v1.NewGetVariablesRequest(jobKey, names...))
	if err != nil {
		return v1.NewInputs(err)
	}
	var vars []*v1.Variable //nolint:prealloc
	for _, v := range variables.Variables {
		vars = append(vars, v)
	}
	return v1.NewInputs(err, vars...)
}

func (g VariableHandler) Input(jobKey, name string) *v1.Input {
	return g.Inputs(jobKey, name).Get(name)
}

func (g VariableHandler) Output(jobKey string, variables ...*models.Variable) error {
	request, err := v1.NewSetVariablesRequest(jobKey, variables...)
	if err != nil {
		return err
	}
	_, err = g.client.SetVariables(context.Background(), request)
	return err
}
