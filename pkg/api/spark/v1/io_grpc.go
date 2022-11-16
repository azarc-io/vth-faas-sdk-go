package spark_v1

import (
	"context"
)

type VariableHandler struct {
	client ManagerServiceClient
}

func NewIOHandler(client ManagerServiceClient) IOHandler {
	return VariableHandler{client}
}

func (g VariableHandler) Inputs(jobKey string, names ...string) *Inputs {
	variables, err := g.client.GetVariables(context.Background(), NewGetVariablesRequest(jobKey, names...))
	if err != nil {
		return NewInputs(err)
	}
	var vars []*Variable //nolint:prealloc
	for _, v := range variables.Variables {
		vars = append(vars, v)
	}
	return NewInputs(err, vars...)
}

func (g VariableHandler) Input(jobKey, name string) *Input {
	return g.Inputs(jobKey, name).Get(name)
}

func (g VariableHandler) Output(jobKey string, variables ...*Var) error {
	request, err := NewSetVariablesRequest(jobKey, variables...)
	if err != nil {
		return err
	}
	_, err = g.client.SetVariables(context.Background(), request)
	return err
}
