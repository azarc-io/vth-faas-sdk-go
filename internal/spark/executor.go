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

//nolint:cyclop
func (c *Chain) runner(ctx sdk_v1.SparkContext, node *Node) sdk_v1.StageError {
	stages := getStagesToResume(node, ctx.LastActiveStage())
	for _, stg := range stages {
		ctx.Log().AddFields(stageLogField, stg.name).AddFields(jobKeyLogField, ctx.JobKey())

		if err := stg.ApplyConditionalExecutionOptions(ctx, stg.name); err != nil {
			er := updateStage(ctx, stg.name, withStageError(err))
			if er != nil {
				ctx.Log().Error(er, "error updating stage status to 'started'")
				return sdk_errors.NewStageError(er)
			}
			continue
		}

		er := updateStage(ctx, stg.name, withStageStatus(sdk_v1.StageStatus_STAGE_STATUS_STARTED))

		if er != nil {
			ctx.Log().Error(er, "error updating stage status to 'started'")
			return sdk_errors.NewStageError(er)
		}

		result, err := stg.cb(context.NewStageContext(ctx))

		if err := c.handleStageError(ctx, node, stg, err); err != nil {
			if err.ErrorType() == sdk_v1.ErrorType_ERROR_TYPE_SKIP {
				continue
			}
			return err
		}

		if err := storeStageResult(ctx, stg, result); err != nil {
			return err
		}

		if err := updateStage(ctx, stg.name, withStageStatus(sdk_v1.StageStatus_STAGE_STATUS_COMPLETED)); err != nil {
			ctx.Log().Error(err, "error setting the stage status to 'completed'")
			return sdk_errors.NewStageError(err)
		}
	}

	if node.complete != nil {
		if er := updateStage(ctx, node.complete.name, withStageStatus(sdk_v1.StageStatus_STAGE_STATUS_STARTED)); er != nil {
			ctx.Log().Error(er, "error setting the completed stage status to 'started'")
			return sdk_errors.NewStageError(er)
		}

		err := node.complete.cb(context.NewCompleteContext(ctx))

		if e := updateStage(ctx, node.complete.name, withStageStatusOrError(sdk_v1.StageStatus_STAGE_STATUS_COMPLETED, err)); e != nil {
			ctx.Log().Error(e, "error setting the completed stage status to 'completed'")
			return sdk_errors.NewStageError(e)
		}
		return err
	}

	return nil
}

//nolint:cyclop
func (c *Chain) handleStageError(ctx sdk_v1.SparkContext, node *Node, stg *stage, err sdk_v1.StageError) sdk_v1.StageError {
	if err == nil {
		return nil
	}
	if e := updateStage(ctx, stg.name, withStageError(err)); e != nil {
		ctx.Log().Error(err, "error updating stage status")
		return sdk_errors.NewStageError(e)
	}
	switch err.ErrorType() {
	case sdk_v1.ErrorType_ERROR_TYPE_FAILED_UNSPECIFIED:
		if node.compensate != nil {
			e := c.runner(ctx.WithoutLastActiveStage(), node.compensate)
			if e != nil {
				return e
			}
		}
		return err
	case sdk_v1.ErrorType_ERROR_TYPE_CANCELLED:
		if node.cancel != nil {
			e := c.runner(ctx.WithoutLastActiveStage(), node.cancel)
			if e != nil {
				return e
			}
		}
		return err
	case sdk_v1.ErrorType_ERROR_TYPE_RETRY:
		return err
	case sdk_v1.ErrorType_ERROR_TYPE_SKIP:
		return err
	case sdk_v1.ErrorType_ERROR_TYPE_FATAL:
		fallthrough
	default:
		ctx.Log().Error(err, "unsupported error type returned from stage '%s'", stg.name)
		return sdk_errors.NewStageError(err, sdk_errors.WithErrorType(sdk_v1.ErrorType_ERROR_TYPE_FATAL))
	}
}

func storeStageResult(ctx sdk_v1.SparkContext, stg *stage, result any) sdk_v1.StageError {
	if result != nil { //nolint:nestif
		req, err := sdk_v1.NewSetStageResultReq(ctx.JobKey(), stg.name, result)
		if err != nil {
			ctx.Log().Error(err, "error creating set stage status request")
			if e := updateStage(ctx, stg.name, withError(err)); e != nil {
				ctx.Log().Error(err, "error updating stage status")
				return sdk_errors.NewStageError(e)
			}
			return sdk_errors.NewStageError(err)
		}
		if err := ctx.StageProgressHandler().SetResult(req); err != nil {
			ctx.Log().Error(err, "error on set stage status request")
			if e := updateStage(ctx, stg.name, withError(err)); e != nil {
				ctx.Log().Error(err, "error updating stage status")
				return sdk_errors.NewStageError(e)
			}
			return sdk_errors.NewStageError(err)
		}
	}
	return nil
}

type updateStageOption = func(stage *sdk_v1.SetStageStatusRequest) *sdk_v1.SetStageStatusRequest

func withStageStatusOrError(status sdk_v1.StageStatus, err sdk_v1.StageError) updateStageOption {
	return func(stage *sdk_v1.SetStageStatusRequest) *sdk_v1.SetStageStatusRequest {
		if err != nil {
			return withStageError(err)(stage)
		}
		return withStageStatus(status)(stage)
	}
}

func withStageStatus(status sdk_v1.StageStatus) updateStageOption {
	return func(stage *sdk_v1.SetStageStatusRequest) *sdk_v1.SetStageStatusRequest {
		stage.Status = status
		return stage
	}
}

func withStageError(err sdk_v1.StageError) updateStageOption {
	return func(stage *sdk_v1.SetStageStatusRequest) *sdk_v1.SetStageStatusRequest {
		if err == nil {
			return stage
		}
		stage.Status = sdk_errors.ErrorTypeToStageStatusMapper(err.ErrorType())
		stage.Err = err.ToErrorMessage()
		return stage
	}
}

func withError(err error) updateStageOption {
	return func(stage *sdk_v1.SetStageStatusRequest) *sdk_v1.SetStageStatusRequest {
		if err == nil {
			return stage
		}
		stage.Status = sdk_v1.StageStatus_STAGE_STATUS_FAILED
		stage.Err = sdk_errors.NewStageError(err).ToErrorMessage()
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

func getStagesToResume(n *Node, lastActiveStage *sdk_v1.LastActiveStage) []*stage {
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
