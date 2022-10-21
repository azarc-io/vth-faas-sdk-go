package context

import (
	"errors"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
)

type Stage struct {
	sdk_v1.Context
	jobContext sdk_v1.SparkContext
}

func (sc Stage) Log() sdk_v1.Logger {
	return sc.jobContext.Log()
}

func NewStageContext(ctx sdk_v1.SparkContext) sdk_v1.StageContext {
	return Stage{jobContext: ctx, Context: ctx}
}

func (sc Stage) Inputs(stage string, names ...string) (*sdk_v1.Variables, error) {
	return sc.jobContext.VariableHandler().Get(sc.JobKey(), stage, names...)
}

func (sc Stage) Input(stage string, name string) (*sdk_v1.Variable, error) {
	vars, err := sc.jobContext.VariableHandler().Get(sc.JobKey(), stage, name)
	if err != nil {
		return nil, err
	}
	if v, ok := vars.Get(name); ok {
		return v, nil
	}
	return nil, errors.New("variable not found") // TODO add context, const error
}
