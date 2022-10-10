package context

import (
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
)

type Compensation struct {
	api.Context
	jobContext   *Job
	stageContext Stage
}

func NewCompensationContext(jobCtx *Job) api.CompensationContext {
	return Compensation{jobContext: jobCtx, Context: jobCtx.metadata}
}

func (c Compensation) Stage(name string, sdf api.StageDefinitionFn, options ...api.StageOption) api.StageChain {
	return c.jobContext.Stage(name, sdf, options...)
}

func (c Compensation) GetVariables(stage string, name ...string) ([]*sdk_v1.Variable, error) {
	return c.jobContext.variableHandler.Get(stage, c.JobKey(), name...)
}

func (c Compensation) SetVariables(stage string, variables ...*sdk_v1.Variable) error {
	return c.jobContext.variableHandler.Set(c.JobKey(), stage, variables...)
}
