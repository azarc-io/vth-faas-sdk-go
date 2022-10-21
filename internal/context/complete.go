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

func (sc Complete) Log() sdk_v1.Logger {
	return sc.jobContext.Log()
}

func (sc Complete) StageResult(name string) (*sdk_v1.StageResult, error) {
	return sc.jobContext.StageProgressHandler().GetResult(sc.JobKey(), name)
}

func (sc Complete) GetVariables(stage string, names ...string) (*sdk_v1.Variables, error) {
	return sc.jobContext.VariableHandler().Get(sc.JobKey(), stage, names...)
}

func (sc Complete) Output(stage string, variables ...*sdk_v1.Variable) error {
	return sc.jobContext.VariableHandler().Set(sc.JobKey(), stage, variables...)
}
