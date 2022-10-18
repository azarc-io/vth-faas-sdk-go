package context

import (
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
)

type Complete struct {
	api.Context
	jobContext api.JobContext
}

func NewCompleteContext(ctx api.JobContext) api.CompleteContext {
	return Complete{jobContext: ctx, Context: ctx}
}

func (sc Complete) GetVariables(stage string, names ...string) (*sdk_v1.Variables, error) {
	return sc.jobContext.VariableHandler().Get(sc.JobKey(), stage, names...)
}

func (sc Complete) SetVariables(stage string, variables ...*sdk_v1.Variable) error {
	return sc.jobContext.VariableHandler().Set(sc.JobKey(), stage, variables...)
}
