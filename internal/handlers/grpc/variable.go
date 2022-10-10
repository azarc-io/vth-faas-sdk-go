package test

import (
	"context"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
)

type GrpcVariableHandler struct {
	client sdk_v1.ManagerServiceClient
}

func NewGrpcVariableHandler() api.VariableHandler {
	return GrpcVariableHandler{}
}

func (g GrpcVariableHandler) Set(req *sdk_v1.SetVariableRequest) error {
	_, err := g.client.SetVariable(context.Background(), req)
	return err
}

func (g GrpcVariableHandler) Get(name, stage, jobKey string) (*sdk_v1.Variable, error) {
	variable, err := g.client.GetVariable(context.Background(), sdk_v1.NewGetVariableRequest(name, stage, jobKey))
	if err != nil {
		return nil, err
	}
	return variable.Variable, nil
}
