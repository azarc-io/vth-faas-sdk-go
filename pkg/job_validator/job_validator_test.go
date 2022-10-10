package job_validator

import (
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	"testing"
)

type ValidExecutor struct{}

var _ api.Job = &ValidExecutor{}

func (i ValidExecutor) Initialize() error {
	return nil
}

func (i ValidExecutor) Execute(ctx api.JobContext) {
	ctx.Stage("stage1", func(stageContext api.StageContext) (any, api.StageError) {
		return nil, nil
	}).Stage("stage2", func(stageContext api.StageContext) (any, api.StageError) {
		return nil, nil
	}).Complete(func(context api.CompletionContext) api.StageError {
		return nil
	}).Compensate(func(compensateContext api.CompensationContext) api.StageError {
		compensateContext.Stage("compensate_stage1", func(c api.StageContext) (any, api.StageError) {
			return nil, nil
		}).Stage("compensate_stage2", func(c api.StageContext) (any, api.StageError) {
			return nil, nil
		}).Canceled(func(cancelContext api.CancelContext) api.StageError {
			cancelContext.Stage("cancel_stage1", func(stageContext api.StageContext) (any, api.StageError) {
				return nil, nil
			}).Compensate(func(context api.CompensationContext) api.StageError {
				context.Stage("another_compensate", func(stageContext api.StageContext) (any, api.StageError) {
					return nil, nil
				})
				return nil
			})
			return nil
		})
		return nil
	})
	return
}

func TestJobExecutorValidatorWithValidExecutor(t *testing.T) {
	executor := ValidExecutor{}
	err := Check(executor)
	if err != nil {
		t.Errorf("job executor is valid and we got this error: %s", err.Error())
	}
}

type InvalidExecutor struct{}

func (i InvalidExecutor) Initialize() error {
	return nil
}

var _ api.Job = &InvalidExecutor{}

func (i InvalidExecutor) Execute(ctx api.JobContext) {
	ctx.Stage("stage1", func(stageContext api.StageContext) (any, api.StageError) {
		return nil, nil
	}).Stage("stage2", func(stageContext api.StageContext) (any, api.StageError) {
		return nil, nil
	}).Complete(func(context api.CompletionContext) api.StageError {
		return nil
	}).Compensate(func(compensateContext api.CompensationContext) api.StageError {
		compensateContext.Stage("compensate_stage1", func(c api.StageContext) (any, api.StageError) {
			return nil, nil
		}).Stage("compensate_stage2", func(c api.StageContext) (any, api.StageError) {
			return nil, nil
		}).Canceled(func(cancelContext api.CancelContext) api.StageError {
			cancelContext.Stage("cancel_stage1", func(stageContext api.StageContext) (any, api.StageError) {
				return nil, nil
			}).Compensate(func(context api.CompensationContext) api.StageError {
				context.Stage("compensate_stage1", func(stageContext api.StageContext) (any, api.StageError) {
					return nil, nil
				}).Canceled(func(context api.CancelContext) api.StageError {
					context.Stage("cancel_stage1", func(context api.StageContext) (any, api.StageError) {
						return nil, nil
					}).Compensate(func(context api.CompensationContext) api.StageError {
						context.Stage("stage1", func(context api.StageContext) (any, api.StageError) {
							return nil, nil
						})
						return nil
					})
					return nil
				})
				return nil
			})
			return nil
		})
		return nil
	})
	return
}

func TestJobExecutorValidatorWithInvalidExecutor(t *testing.T) {
	expectedErrorMessage := "invalid stage chain. stage names must be unique. the following stage names have more than one occurrence: cancel_stage1, compensate_stage1, stage1"
	executor := InvalidExecutor{}
	err := Check(executor)
	if err == nil {
		t.Error("job executor is invalid and an error was expected.")
	}
	if err.Error() != expectedErrorMessage {
		t.Errorf(`invalid job executor error message not matching with the expected
want: '%s'
got:  '%s'`, err.Error(), expectedErrorMessage)
	}
}
