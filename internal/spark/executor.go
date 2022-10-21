package spark

import (
	"github.com/azarc-io/vth-faas-sdk-go/internal/context"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	sdk_errors "github.com/azarc-io/vth-faas-sdk-go/pkg/errors"
)

const (
	stageLogField  = "stage"
	jobKeyLogField = "job_key"
)

func (c *Chain) Execute(ctx api.SparkContext) api.StageError {
	return c.runner(ctx, c.rootNode)
}

func (c *Chain) runner(ctx api.SparkContext, node *node) api.StageError {
	for _, stg := range node.stages {
		ctx.Log().AddFields(stageLogField, stg.name).AddFields(jobKeyLogField, ctx.JobKey())

		if err := stg.ApplyStageOptionsParams(ctx, stg.name); err != nil {
			return err
		}

		er := updateStage(ctx, stg.name, withStageStatus(sdk_v1.StageStatus_StageStarted))
		if er != nil {
			ctx.Log().Error(er, "error trying to update a stage status")
			return sdk_errors.NewStageError(er)
		}

		result, err := stg.cb(context.NewStageContext(ctx))

		if err != nil {
			if e := updateStage(ctx, stg.name, withStageError(err)); e != nil {
				ctx.Log().Error(err, "error trying to update a stage status")
				return sdk_errors.NewStageError(e)
			}
			switch err.ErrorType() {
			case sdk_v1.ErrorType_Failed:
				if node.compensate != nil {
					err := c.runner(ctx, node.compensate)
					return err
				}
				return err
			case sdk_v1.ErrorType_Canceled:
				if node.cancel != nil {
					return c.runner(ctx, node.cancel)
				}
				return err
			case sdk_v1.ErrorType_Retry:
				return err
			case sdk_v1.ErrorType_Skip:
				continue // =)
			default:
				// TODO log error unsupported error
				return err
			}
		}
		if result != nil {
			req, err := sdk_v1.NewSetStageResultReq(ctx.JobKey(), stg.name, result)
			if err != nil {
				// TODO log error -->> retry forever?
				return sdk_errors.NewStageError(err)
			}
			if err := ctx.StageProgressHandler().SetResult(req); err != nil {
				// TODO log error -->> retry forever?
				return sdk_errors.NewStageError(err)
			}
		}
		if err := updateStage(ctx, stg.name, withStageStatus(sdk_v1.StageStatus_StageCompleted)); err != nil {
			// TODO log error -->> retry forever?
			return sdk_errors.NewStageError(err)
		}
	}
	if node.complete != nil {
		node.complete.cb(context.NewCompleteContext(ctx))
	}
	return nil
}

type updateStageOption = func(stage *sdk_v1.SetStageStatusRequest) *sdk_v1.SetStageStatusRequest

func withStageStatus(status sdk_v1.StageStatus) updateStageOption {
	return func(stage *sdk_v1.SetStageStatusRequest) *sdk_v1.SetStageStatusRequest {
		stage.Status = status
		return stage
	}
}

func withStageError(err api.StageError) updateStageOption {
	return func(stage *sdk_v1.SetStageStatusRequest) *sdk_v1.SetStageStatusRequest {
		stage.Status = sdk_errors.ErrorTypeToStageStatusMapper[err.ErrorType()]
		stage.Err = err.ToErrorMessage()
		return stage
	}
}

func updateStage(ctx api.SparkContext, name string, opts ...updateStageOption) error {
	req := &sdk_v1.SetStageStatusRequest{JobKey: ctx.JobKey(), Name: name}
	for _, opt := range opts {
		req = opt(req)
	}
	return ctx.StageProgressHandler().Set(req)
}
