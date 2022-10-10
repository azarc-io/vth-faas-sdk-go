package context

import (
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
)

type Cancellation struct {
	api.Context
	jobContext   *Job
	stageContext Stage
}

func NewCancellationContext(jobCtx *Job) api.CancelContext {
	return Cancellation{jobContext: jobCtx, Context: jobCtx.metadata}
}

func (c Cancellation) Stage(name string, sdf api.StageDefinitionFn, options ...api.StageOption) api.StageChain {
	return c.jobContext.Stage(name, sdf, options...)
}
