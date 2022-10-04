package context

import sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"

type Stage struct {
	jobContext Job
}

func NewStageContext(jobCtx Job) Stage {
	return Stage{jobCtx}
}

func (sc Stage) GetVariable(name string) (*sdk_v1.Variable, error) {
	return sc.jobContext.variableHandler.Get(name)
}
