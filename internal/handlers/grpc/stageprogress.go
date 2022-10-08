package test

import (
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	sdk_errors "github.com/azarc-io/vth-faas-sdk-go/pkg/errors"
)

// TODO implement the DEFAULT GRPC handler
// this must also trigger the spark execution
type inMemoryStageProgressHandler struct {
	stages map[string]*sdk_v1.Stage
}

func NewMockStageProgressHandler(stages ...*sdk_v1.Stage) api.StageProgressHandler {
	handler := inMemoryStageProgressHandler{map[string]*sdk_v1.Stage{}}
	if stages == nil {
		return handler
	}
	for _, stage := range stages {
		handler.stages[stage.Name] = stage
	}
	return handler
}

func (i inMemoryStageProgressHandler) Get(name string) (*sdk_v1.Stage, error) {
	if variable, ok := i.stages[name]; ok {
		return variable, nil
	}
	return nil, sdk_errors.VariableNotFound
}

func (i inMemoryStageProgressHandler) Set(stage *sdk_v1.Stage) error {
	i.stages[stage.Name] = stage
	return nil
}
