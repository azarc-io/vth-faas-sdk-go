package context

import sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"

type Completion struct {
	jobContext Job
}

func NewCompleteContext(jobCtx Job) Completion {
	return Completion{jobCtx}
}

func (c Completion) GetStage(name string) (*sdk_v1.Stage, error) {
	return c.jobContext.stageProgressHandler.Get(name)
}

func (c Completion) SetVariable(variable *sdk_v1.Variable) error {
	return c.jobContext.variableHandler.Set(variable)
}

func (c Completion) GetStageResult(stage *sdk_v1.Stage) (*sdk_v1.StageResult, error) {
	return c.jobContext.stageProgressHandler.GetResult(stage)
}
