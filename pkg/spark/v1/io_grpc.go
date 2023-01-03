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

func (g variableHandler) Inputs(ctx SparkContext, names ...string) Inputs {
	request := newGetVariablesRequest(ctx, names...)
	variables, err := g.client.GetInputs(context.Background(), request)
	if err != nil {
		return newInputs(err, nil)
	}

	return newInputs(err, variables.Variables)
}

func (g variableHandler) Input(ctx SparkContext, name string) Input {
	return g.Inputs(ctx, name).Get(name)
}

func (g variableHandler) Output(ctx SparkContext, variables ...*Var) error {
	request, err := newSetVariablesRequest(ctx, variables...)
	if err != nil {
		return err
	}
	_, err = g.client.SetOutputs(context.Background(), request)
	return err
}
