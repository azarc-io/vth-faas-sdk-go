package context

import (
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
)

type Completion struct {
	api.Context
	jobContext *Job
}

func NewCompleteContext(jobCtx *Job) api.CompletionContext {
	return Completion{jobContext: jobCtx, Context: jobCtx.metadata}
}

func (c Completion) GetStageResult(stageName string) (*sdk_v1.StageResult, error) {
	return c.jobContext.stageProgressHandler.GetResult(c.JobKey(), stageName)
}

func (c Completion) SetVariables(stage string, variables ...*sdk_v1.Variable) error {
	return c.jobContext.variableHandler.Set(c.JobKey(), stage, variables...)
}
