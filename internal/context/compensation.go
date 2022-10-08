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

func (c Compensation) Stage(name string, sdf api.StageDefinitionFn) api.StageChain {
	return c.jobContext.Stage(name, sdf)
}

func (c Compensation) WithStageStatus(names []string, value any) bool {
	//TODO implement me // TODO <<<<<<<<<<<<< don't forget this
	panic("implement me")
}

func (c Compensation) GetVariable(s string) (*sdk_v1.Variable, error) {
	return c.jobContext.variableHandler.Get(s)
}

func (c Compensation) SetVariable(name string, value any, mimeType string) error {
	return c.jobContext.variableHandler.Set(sdk_v1.NewVariable(name, mimeType, value))
}
