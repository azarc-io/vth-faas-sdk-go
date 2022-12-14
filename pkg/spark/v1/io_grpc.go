package spark_v1

import (
	"context"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/internal/gen/azarc/sdk/spark/v1"
)

type variableHandler struct {
	client sparkv1.ManagerServiceClient
}

func newGrpcIOHandler(client sparkv1.ManagerServiceClient) IOHandler {
	return variableHandler{client}
}

func (g variableHandler) Inputs(jobKey string, names ...string) Inputs {
	variables, err := g.client.GetInputs(context.Background(), newGetVariablesRequest(jobKey, names...))
	if err != nil {
		return newInputs(err, nil)
	}

	return newInputs(err, variables.Variables)
}

func (g variableHandler) Input(jobKey, name string) Input {
	return g.Inputs(jobKey, name).Get(name)
}

func (g variableHandler) Output(jobKey string, variables ...*Var) error {
	request, err := newSetVariablesRequest(jobKey, variables...)
	if err != nil {
		return err
	}
	_, err = g.client.SetOutputs(context.Background(), request)
	return err
}
