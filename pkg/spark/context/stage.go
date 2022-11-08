package context

import (
	v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1/models"
)

type Stage struct {
	v1.Context
	jobContext v1.SparkContext
}

func NewCompleteContext(ctx v1.SparkContext) v1.CompleteContext {
	return Stage{jobContext: ctx, Context: ctx}
}

func NewStageContext(ctx v1.SparkContext) v1.StageContext {
	return Stage{jobContext: ctx, Context: ctx}
}

func (sc Stage) Inputs(names ...string) *v1.Inputs {
	return sc.jobContext.IOHandler().Inputs(sc.JobKey(), names...)
}

func (sc Stage) Input(name string) *v1.Input {
	return sc.jobContext.IOHandler().Input(sc.JobKey(), name)
}

func (sc Stage) StageResult(name string) *v1.Result {
	return sc.jobContext.StageProgressHandler().GetResult(sc.JobKey(), name)
}

func (sc Stage) Output(variables ...*models.Variable) error {
	return sc.jobContext.IOHandler().Output(sc.JobKey(), variables...)
}

func (sc Stage) Log() v1.Logger {
	return sc.jobContext.Log()
}
