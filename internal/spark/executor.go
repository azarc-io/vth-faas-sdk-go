package spark

import (
	"github.com/azarc-io/vth-faas-sdk-go/internal/context"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	sdk_errors "github.com/azarc-io/vth-faas-sdk-go/pkg/errors"
)

const (
	stageLogField  = "stage"
	jobKeyLogField = "job_key"
)

func (c *Chain) Execute(ctx sdk_v1.SparkContext) sdk_v1.StageError {
	n, err := c.getNodeToResume(ctx.LastActiveStage())
	if err != nil {
		return sdk_errors.NewStageError(err)
	}
	return c.runner(ctx, n)
}

func (c *Chain) runner(ctx sdk_v1.SparkContext, node *node) sdk_v1.StageError {
	stages := getStagesToResume(node, ctx.LastActiveStage())
	for _, stg := range stages {
		ctx.Log().AddFields(stageLogField, stg.name).AddFields(jobKeyLogField, ctx.JobKey())

		if err := stg.ApplyStageOptionsParams(ctx, stg.name); err != nil {
			return err
		}

		er := updateStage(ctx, stg.name, withStageStatus(sdk_v1.StageStatus_StageStarted))
		if er != nil {
			ctx.Log().Error(er, "error updating stage status to 'started'")
			return sdk_errors.NewStageError(er)
		}

		result, err := stg.cb(context.NewStageContext(ctx))

		if err != nil {
			if e := updateStage(ctx, stg.name, withStageError(err)); e != nil {
				ctx.Log().Error(err, "error updating stage status")
				return sdk_errors.NewStageError(e)
			}
			switch err.ErrorType() {
			case sdk_v1.ErrorType_Failed:
				if node.compensate != nil {
					return c.runner(ctx, node.compensate)
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
				continue
			default:
				ctx.Log().Error(err, "unsupported error type returned from stage '%s'", stg.name)
				return err
			}
		}
		if result != nil {
			req, err := sdk_v1.NewSetStageResultReq(ctx.JobKey(), stg.name, result)
			if err != nil {
				ctx.Log().Error(err, "error creating set stage status request")
				return sdk_errors.NewStageError(err)
			}
			if err := ctx.StageProgressHandler().SetResult(req); err != nil {
				ctx.Log().Error(err, "error on set stage status request")
				return sdk_errors.NewStageError(err)
			}
		}
		if err := updateStage(ctx, stg.name, withStageStatus(sdk_v1.StageStatus_StageCompleted)); err != nil {
			ctx.Log().Error(err, "error setting the stage status to 'completed'")
			return sdk_errors.NewStageError(err)
		}
	}
	if node.complete != nil {
		return node.complete.cb(context.NewCompleteContext(ctx))
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

func withStageError(err sdk_v1.StageError) updateStageOption {
	return func(stage *sdk_v1.SetStageStatusRequest) *sdk_v1.SetStageStatusRequest {
		stage.Status = sdk_errors.ErrorTypeToStageStatusMapper[err.ErrorType()]
		stage.Err = err.ToErrorMessage()
		return stage
	}
}

func updateStage(ctx sdk_v1.SparkContext, name string, opts ...updateStageOption) error {
	req := &sdk_v1.SetStageStatusRequest{JobKey: ctx.JobKey(), Name: name}
	for _, opt := range opts {
		req = opt(req)
	}
	return ctx.StageProgressHandler().Set(req)
}

func getStagesToResume(n *node, lastActiveStage *sdk_v1.LastActiveStage) []*stage {
	if lastActiveStage == nil {
		return n.stages
	}
	var stages []*stage
	for idx, stg := range n.stages {
		if stg.name == lastActiveStage.Name {
			stages = append(stages, n.stages[idx:]...)
		}
	}
	return stages
}
