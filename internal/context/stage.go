package context

import (
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
)

type Stage struct {
	api.Context
	jobContext api.JobContext
}

func NewStageContext(ctx api.JobContext) api.StageContext {
	return Stage{jobContext: ctx, Context: ctx}
}

func (sc Stage) GetVariables(stage string, names ...string) (*sdk_v1.Variables, error) {
	return sc.jobContext.VariableHandler().Get(sc.JobKey(), stage, names...)
}

func (sc Stage) SetVariables(stage string, variables ...*sdk_v1.Variable) error {
	return sc.jobContext.VariableHandler().Set(sc.JobKey(), stage, variables...)
}
