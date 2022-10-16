package job

import (
	"github.com/azarc-io/vth-faas-sdk-go/internal/context"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	sdk_errors "github.com/azarc-io/vth-faas-sdk-go/pkg/errors"
)

func (c *Chain) Execute(ctx api.JobContext) api.StageError {
	return c.runner(ctx, c.rootNode)
}

func (c *Chain) runner(ctx api.JobContext, node *node) api.StageError {
	for _, stg := range node.stages {
		er := updateStage(ctx, stg.name, withStageStatus(sdk_v1.StageStatus_StageStarted))
		if er != nil {
			// TODO log error -->> retry forever?
			return sdk_errors.NewStageError(er)
		}
		result, err := stg.cb(context.NewStageContext(ctx))
		if err != nil {
			if e := updateStage(ctx, stg.name, withStageError(err)); e != nil {
				// TODO log error -->> retry forever?
				return sdk_errors.NewStageError(e)
			}
			switch err.ErrorType() {
			case sdk_v1.ErrorType_Failed:
				if node.compensate != nil {
					// TODO update stage to compensate starting
					c.runner(ctx, node.compensate)
					// TODO update stage to compensate done / errored
				}
				return err
			case sdk_v1.ErrorType_Canceled:
				if node.cancel != nil {
					c.runner(ctx, node.cancel)
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

func updateStage(ctx api.JobContext, name string, opts ...updateStageOption) error {
	req := &sdk_v1.SetStageStatusRequest{JobKey: ctx.JobKey(), Name: name}
	for _, opt := range opts {
		req = opt(req)
	}
	return ctx.StageProgressHandler().Set(req)
}
