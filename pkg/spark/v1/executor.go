package spark_v1

import (
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/internal/gen/azarc/sdk/spark/v1"
	"github.com/rs/zerolog/log"
)

const (
	stageLogField  = "stage"
	jobKeyLogField = "job_key"
)

func (c *chain) execute(ctx SparkContext) StageError {
	n, err := c.getNodeToResume(ctx.LastActiveStage())
	if err != nil {
		log.Error().Err(err).Msgf("failed to resolve node to resume from")
		return NewStageError(err)
	}
	if err := c.runner(ctx, n); err != nil {
		return err
	}

	if err := finishJob(ctx); err != nil {
		return NewStageError(err, WithErrorCode(500), WithMetadata(map[string]string{
			"reason": "could not finalize the job",
		}))
	}

	return nil
}

//nolint:cyclop
func (c *chain) runner(ctx SparkContext, node *node) StageError {
	stages := getStagesToResume(node, ctx.LastActiveStage())
	var names []string
	for _, s := range stages {
		names = append(names, s.name)
	}
	log.Info().Msgf("running stages: %v", names)
	for _, stg := range stages {
		select {
		case <-ctx.Ctx().Done():
			ctx.Log().Debug("exited due to context completion")
			return nil
		default:
			ctx.Log().AddFields(stageLogField, stg.name).AddFields(jobKeyLogField, ctx.JobKey())

			if err := stg.ApplyConditionalExecutionOptions(ctx, stg.name); err != nil {
				er := updateStage(ctx, stg.name, withStageError(err))
				if er != nil {
					ctx.Log().Error(er, "error updating stage status to 'started'")
					return NewStageError(er)
				}
				continue
			}

			er := updateStage(ctx, stg.name, withStageStatus(sparkv1.StageStatus_STAGE_STARTED))

			if er != nil {
				ctx.Log().Error(er, "error updating stage status to 'started'")
				return NewStageError(er)
			}

			var result any
			var stageErr StageError

			// stage execution is delegated in which case call the delegate
			// instead and expect that it will invoke the stage and return a result, error
			if ctx.delegateStage() != nil {
				result, stageErr = ctx.delegateStage()(NewStageContext(ctx, stg.name), stg.cb)
			} else {
				result, stageErr = stg.cb(NewStageContext(ctx, stg.name))
			}

			if err := c.handleStageError(ctx, node, stg, stageErr); err != nil {
				if err.ErrorType() == sparkv1.ErrorType_ERROR_TYPE_SKIP {
					continue
				}
				ctx.Log().Error(err, "could not sync stage error")
				return err
			}

			if err := storeStageResult(ctx, stg, result); err != nil {
				ctx.Log().Error(err, "could not sync stage result")
				return err
			}

			if err := updateStage(ctx, stg.name, withStageStatus(sparkv1.StageStatus_STAGE_COMPLETED)); err != nil {
				ctx.Log().Error(err, "error setting the stage status to 'completed'")
				return NewStageError(err)
			}
		}
	}

	select {
	case <-ctx.Ctx().Done():
		return nil
	default:
	}

	if node.complete != nil {
		if er := updateStage(ctx, node.complete.name, withStageStatus(sparkv1.StageStatus_STAGE_STARTED)); er != nil {
			ctx.Log().Error(er, "error setting the completed stage status to 'started'")
			return NewStageError(er)
		}

		var stageErr StageError

		if ctx.delegateComplete() != nil {
			stageErr = ctx.delegateComplete()(NewCompleteContext(ctx, node.complete.name), node.complete.cb)
		} else {
			stageErr = node.complete.cb(NewCompleteContext(ctx, node.complete.name))
		}

		if e := updateStage(ctx, node.complete.name, withStageStatusOrError(sparkv1.StageStatus_STAGE_COMPLETED, stageErr)); e != nil {
			ctx.Log().Error(e, "error setting the completed stage status to 'completed'")
			return NewStageError(e)
		}
		return stageErr
	}

	return nil
}

//nolint:cyclop
func (c *chain) handleStageError(ctx SparkContext, node *node, stg *stage, err StageError) StageError {
	if err == nil {
		return nil
	}

	log.Error().Err(err).Msgf("caught error while executing stage: %s", stg.name)

	if e := updateStage(ctx, stg.name, withStageError(err)); e != nil {
		ctx.Log().Error(err, "error updating stage status during stage error")
		return NewStageError(e)
	}

	switch err.ErrorType() {
	case sparkv1.ErrorType_ERROR_TYPE_FAILED_UNSPECIFIED:
		if node.compensate != nil {
			e := c.runner(ctx.WithoutLastActiveStage(), node.compensate)
			if e != nil {
				return e
			}
		}
		return err
	case sparkv1.ErrorType_ERROR_TYPE_CANCELLED:
		if node.cancel != nil {
			e := c.runner(ctx.WithoutLastActiveStage(), node.cancel)
			if e != nil {
				return e
			}
		}
		return err
	case sparkv1.ErrorType_ERROR_TYPE_RETRY:
		return err
	case sparkv1.ErrorType_ERROR_TYPE_SKIP:
		return err
	case sparkv1.ErrorType_ERROR_TYPE_FATAL:
		fallthrough
	default:
		ctx.Log().Error(err, "unsupported error type returned from stage '%s'", stg.name)
		return NewStageError(err, withErrorType(sparkv1.ErrorType_ERROR_TYPE_FATAL))
	}
}

func storeStageResult(ctx SparkContext, stg *stage, result any) StageError {
	if result != nil { //nolint:nestif
		req, err := newSetStageResultReq(ctx, stg.name, result)
		if err != nil {
			ctx.Log().Error(err, "error creating set stage status request")
			if e := updateStage(ctx, stg.name, withError(err)); e != nil {
				ctx.Log().Error(err, "failed to generate stage status request")
				return NewStageError(e)
			}
			return NewStageError(err)
		}
		if err := ctx.StageProgressHandler().SetResult(ctx, req); err != nil {
			ctx.Log().Error(err, "error on set stage status request")
			if e := updateStage(ctx, stg.name, withError(err)); e != nil {
				ctx.Log().Error(err, "error updating stage status")
				return NewStageError(e)
			}
			return NewStageError(err)
		}
	}
	return nil
}

type updateStageOption = func(stage *sparkv1.SetStageStatusRequest) *sparkv1.SetStageStatusRequest

func withStageStatusOrError(status sparkv1.StageStatus, err StageError) updateStageOption {
	return func(stage *sparkv1.SetStageStatusRequest) *sparkv1.SetStageStatusRequest {
		if err != nil {
			return withStageError(err)(stage)
		}
		return withStageStatus(status)(stage)
	}
}

func withStageStatus(status sparkv1.StageStatus) updateStageOption {
	return func(stage *sparkv1.SetStageStatusRequest) *sparkv1.SetStageStatusRequest {
		stage.Status = status
		return stage
	}
}

func withStageError(err StageError) updateStageOption {
	return func(stage *sparkv1.SetStageStatusRequest) *sparkv1.SetStageStatusRequest {
		if err == nil {
			return stage
		}
		stage.Status = errorTypeToStageStatus(err.ErrorType())
		stage.Err = err.ToErrorMessage()
		return stage
	}
}

func withError(err error) updateStageOption {
	return func(stage *sparkv1.SetStageStatusRequest) *sparkv1.SetStageStatusRequest {
		if err == nil {
			return stage
		}
		stage.Status = sparkv1.StageStatus_STAGE_FAILED
		stage.Err = NewStageError(err).ToErrorMessage()
		return stage
	}
}

func updateStage(ctx SparkContext, name string, opts ...updateStageOption) error {
	req := &sparkv1.SetStageStatusRequest{Key: ctx.JobKey(), Name: name, Metadata: &sparkv1.RequestMetadata{
		Metadata: ctx.RequestMetadata(),
	}}
	for _, opt := range opts {
		req = opt(req)
	}
	return ctx.StageProgressHandler().Set(ctx, req)
}

func finishJob(ctx SparkContext) error {
	req := &sparkv1.FinishJobRequest{Key: ctx.JobKey(), Metadata: &sparkv1.RequestMetadata{
		Metadata: ctx.RequestMetadata(),
	}}

	return ctx.StageProgressHandler().FinishJob(ctx, req)
}

func getStagesToResume(n *node, lastActiveStage *sparkv1.LastActiveStage) []*stage {
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
