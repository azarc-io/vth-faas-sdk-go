package spark_v1

import (
	"context"
)

type variableHandler struct {
	client ManagerServiceClient
}

func newGrpcIOHandler(client ManagerServiceClient) IOHandler {
	return variableHandler{client}
}

func (g variableHandler) Inputs(jobKey string, names ...string) Inputs {
	variables, err := g.client.GetVariables(context.Background(), newGetVariablesRequest(jobKey, names...))
	if err != nil {
		return newInputs(err)
	}
	var vars []*Variable //nolint:prealloc
	for _, v := range variables.Variables {
		vars = append(vars, v)
	}
	return newInputs(err, vars...)
}

func (g variableHandler) Input(jobKey, name string) Input {
	return g.Inputs(jobKey, name).Get(name)
}

func (g variableHandler) Output(jobKey string, variables ...*Var) error {
	request, err := newSetVariablesRequest(jobKey, variables...)
	if err != nil {
		return err
	}
	_, err = g.client.SetVariables(context.Background(), request)
	return err
}
