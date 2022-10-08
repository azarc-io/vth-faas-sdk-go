package context

import (
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
)

type Stage struct {
	api.Context
	jobContext *Job
}

func NewStageContext(jobCtx *Job) Stage {
	return Stage{jobContext: jobCtx, Context: jobCtx.metadata}
}

func (sc Stage) GetVariable(name string) (*sdk_v1.Variable, error) {
	return sc.jobContext.variableHandler.Get(name)
}
