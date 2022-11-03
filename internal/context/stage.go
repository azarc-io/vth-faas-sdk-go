package context

import (
	"github.com/azarc-io/vth-faas-sdk-go/internal/handlers"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
)

type Stage struct {
	sdk_v1.Context
	jobContext sdk_v1.SparkContext
}

func NewCompleteContext(ctx sdk_v1.SparkContext) sdk_v1.CompleteContext {
	return Stage{jobContext: ctx, Context: ctx}
}

func NewStageContext(ctx sdk_v1.SparkContext) sdk_v1.StageContext {
	return Stage{jobContext: ctx, Context: ctx}
}

func (sc Stage) Inputs(names ...string) *sdk_v1.Inputs {
	return sc.jobContext.IOHandler().Inputs(sc.JobKey(), names...)
}

func (sc Stage) Input(name string) *sdk_v1.Input {
	return sc.jobContext.IOHandler().Input(sc.JobKey(), name)
}

func (sc Stage) StageResult(name string) *sdk_v1.Result {
	return sc.jobContext.StageProgressHandler().GetResult(sc.JobKey(), name)
}

func (sc Stage) Output(variables ...*handlers.Variable) error {
	return sc.jobContext.IOHandler().Output(sc.JobKey(), variables...)
}

func (sc Stage) Log() sdk_v1.Logger {
	return sc.jobContext.Log()
}
