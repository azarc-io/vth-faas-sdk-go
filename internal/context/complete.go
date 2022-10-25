package context

import (
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
)

type Complete struct {
	sdk_v1.Context
	jobContext sdk_v1.SparkContext
}

func NewCompleteContext(ctx sdk_v1.SparkContext) sdk_v1.CompleteContext {
	return Complete{jobContext: ctx, Context: ctx}
}

func (sc Complete) Inputs(names ...string) *sdk_v1.Inputs {
	return sc.jobContext.VariableHandler().Get(sc.JobKey(), names...)
}

func (sc Complete) Input(name string) *sdk_v1.Input {
	return sc.jobContext.VariableHandler().Get(sc.JobKey(), name).Get(name)
}

func (sc Complete) Output(variables ...*sdk_v1.Variable) error {
	return sc.jobContext.VariableHandler().Set(sc.JobKey(), variables...)
}

func (sc Complete) StageResult(name string) (*sdk_v1.StageResult, error) {
	return sc.jobContext.StageProgressHandler().GetResult(sc.JobKey(), name)
}

func (sc Complete) Log() sdk_v1.Logger {
	return sc.jobContext.Log()
}
