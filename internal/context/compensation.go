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

func (c Compensation) GetVariable(name, stage string) (*sdk_v1.Variable, error) {
	return c.jobContext.variableHandler.Get(name, stage, c.JobKey())
}

func (c Compensation) SetVariable(variable *sdk_v1.SetVariableRequest) error {
	return c.jobContext.variableHandler.Set(variable)
}
