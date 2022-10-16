package main

import (
	ctx "context"
	"github.com/azarc-io/vth-faas-sdk-go/internal/context"
	"github.com/azarc-io/vth-faas-sdk-go/internal/handlers/test/inmemory"
	"github.com/azarc-io/vth-faas-sdk-go/internal/worker"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/config"
	sdk_errors "github.com/azarc-io/vth-faas-sdk-go/pkg/errors"
	"testing"
)

type SampleJob struct {
}

func (s SampleJob) Initialize() error {
	return nil
}

func (s SampleJob) Execute(ctx api.JobContext) {
	ctx.Stage("stage1", func(stageCtx api.StageContext) (any, api.StageError) {
		variables, err := stageCtx.GetVariables("stage1", "var1", "var2", "var3") // <--
		if err != nil {
			return nil, sdk_errors.NewStageError(err, sdk_errors.WithRetry(10, 1000))
		}
		if v, ok := variables.Get("var2"); ok {
			var value map[string]any
			err = v.Bind(&value)
			if err != nil {
				return nil, sdk_errors.NewStageError(err)
			}
			return value, nil
		}
		// >> do your business logic here <<
		return nil, nil
	}).Stage("stage2", func(stageCtx api.StageContext) (any, api.StageError) {
		// >> do your business logic here <<
		return nil, nil
	}).Stage("stage3", func(stageCtx api.StageContext) (any, api.StageError) {
		// >> do your business logic here <<
		return nil, nil
	}).Complete(func(completeCtx api.CompletionContext) api.StageError {
		result, err := completeCtx.GetStageResult("stage1") // <--
		if err != nil {
			return sdk_errors.NewStageError(err, sdk_errors.WithErrorType(sdk_v1.ErrorType_Failed))
		}
		value := map[string]any{}
		_ = result.Bind(&value)

		newVariable, _ := sdk_v1.NewVariable("var4", "application/json", "value")
		_ = completeCtx.SetVariables("stage3", newVariable) // <--
		// >> do your business logic here <<
		return nil
	}).Compensate(func(compensateCtx api.CompensationContext) api.StageError {
		_, _ = compensateCtx.GetVariables("stage1", "var1", "var2") // < --

		newVariable, _ := sdk_v1.NewVariable("var4", "application/json", "value")
		_ = compensateCtx.SetVariables("stage3", newVariable) // <--

		compensateCtx.Stage("compensate-stage-1", func(stageCtx api.StageContext) (any, api.StageError) {
			// >> do your business logic here <<
			return nil, nil
		}, context.WithStageStatus("stage1", sdk_v1.StageStatus_StageFailed)).
			Stage("compensate-stage-2", func(context api.StageContext) (any, api.StageError) {
				// >> do your business logic here <<
				return nil, nil
			}).Compensate(func(compensationContext api.CompensationContext) api.StageError {
			compensationContext.Stage("another-compensate-stage", func(stageContext api.StageContext) (any, api.StageError) {
				// >> do your business logic here <<
				return nil, nil
			})
			return nil
		})
		return nil
	}).Canceled(func(cancelCtx api.CancelContext) api.StageError {
		cancelCtx.Stage("cancel", func(stageContext api.StageContext) (any, api.StageError) {
			return nil, nil
		}, context.WithStageStatus("compensate-stage-2", sdk_v1.StageStatus_StageFailed)).
			Canceled(func(cancelContext api.CancelContext) api.StageError {
				return nil
			},
			)
		return nil
	})
}

func TestCreatingJob(t *testing.T) {
	cfg, err := config.NewMock(map[string]string{"APP_ENVIRONMENT": "test", "AGENT_SERVER_PORT": "0", "MANAGER_SERVER_PORT": "0"})
	if err != nil {
		t.Error(err)
	}

	job := SampleJob{}

	stageProgressHandler := inmemory.NewStageProgressHandler(t,
		sdk_v1.NewSetStageStatusReq("jobKey", "stage1", sdk_v1.StageStatus_StagePending),
		sdk_v1.NewSetStageStatusReq("jobKey", "stage2", sdk_v1.StageStatus_StagePending),
		sdk_v1.NewSetStageStatusReq("jobKey", "stage3", sdk_v1.StageStatus_StagePending),
		sdk_v1.NewSetStageStatusReq("jobKey", "compensate-stage-2", sdk_v1.StageStatus_StagePending),
	)
	var1, _ := sdk_v1.NewVariable("var1", "application/json", map[string]any{"key": "value"})
	var2, _ := sdk_v1.NewVariable("var2", "application/json", map[string]any{"key": "value"})
	var3, _ := sdk_v1.NewVariable("var3", "application/json", map[string]any{"key": "value"})
	variablesHandler := inmemory.NewVariableHandler(t,
		sdk_v1.NewSetVariablesRequest("jobKey", "stage1", var1, var2, var3),
	)

	jobWorker, err := worker.NewJobWorker(cfg, job,
		worker.WithStageProgressHandler(stageProgressHandler),
		worker.WithVariableHandler(variablesHandler))
	if err != nil {
		t.Error(err)
	}
	err = jobWorker.Run(context.NewJobMetadata(ctx.Background(), "jobKey", "correlationId", "transactionId", nil))
}
