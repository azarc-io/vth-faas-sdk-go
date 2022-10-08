package context

import (
	"context"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	sdk_errors "github.com/azarc-io/vth-faas-sdk-go/pkg/errors"
)

type JobMetadata struct {
	jobKey        string
	correlationId string
	transactionId string
	payload       any
}

func NewJobMetadata(jobKey, correlationId, transactionId string, payload any) JobMetadata {
	return JobMetadata{jobKey, correlationId, transactionId, payload}
}

func (j JobMetadata) JobKey() string {
	return j.jobKey
}

func (j JobMetadata) CorrelationID() string {
	return j.correlationId
}

func (j JobMetadata) TransactionID() string {
	return j.transactionId
}

type Job struct {
	ctx                  context.Context
	metadata             JobMetadata
	stageProgressHandler api.StageProgressHandler
	variableHandler      api.VariableHandler
	stageErr             api.StageError
}

func NewJobContext(metadata JobMetadata, sph api.StageProgressHandler, vh api.VariableHandler) api.JobContext {
	return &Job{metadata: metadata, stageProgressHandler: sph, variableHandler: vh}
}

func (j *Job) Stage(name string, stageDefinitionFn api.StageDefinitionFn) api.StageChain {
	stage, err := j.getStage(name)
	if err != nil {
		return j.handleStageError(err)
	}
	switch *stage {
	case sdk_v1.StageStatus_StagePending:
		err = j.updateStage(j.metadata.jobKey, name, withStageStatus(sdk_v1.StageStatus_StageStarted))
		if err != nil {
			return j.handleStageError(err)
		}

		stageContext := NewStageContext(j)
		result, stageErr := stageDefinitionFn(stageContext)

		if stageErr != nil {
			err = j.updateStage(j.metadata.jobKey, name, withStageError(stageErr))
			if err != nil {
				return j.handleStageError(err)
			}
			return j.handleStageError(stageErr)
		}
		err = j.updateStage(j.metadata.jobKey, name, withStageStatus(sdk_v1.StageStatus_StageCompleted))
		if err != nil {
			return j.handleStageError(err)
		}
		if result != nil {
			err = j.stageProgressHandler.SetResult(sdk_v1.NewSetStageResultReq(j.metadata.jobKey, name, result))
			if err != nil {
				return j.handleStageError(err)
			}
		}
		// TODO log.debug success
		return j
	default:
		// StageFailed = 3;
		// StageSkipped = 4;
		// StageCanceled = 5;
		// TODO log.info stageName stageStatus
		return j
	}
}

func (j *Job) Complete(completionDefinitionFn api.CompletionDefinitionFn) api.CompleteChain {
	if j.stageErr != nil {
		// TODO log.info can't execute job's complete stage because of a previous stage error
		return j
	}
	err := j.stageProgressHandler.SetJobStatus(sdk_v1.NewSetJobStatusReq(j.metadata.jobKey, sdk_v1.JobStatus_JobCompletionStarted))
	if err != nil {
		// TODO log.error
		return j
	}

	completionCtx := NewCompleteContext(j)
	err = completionDefinitionFn(completionCtx)

	if err != nil {
		// TODO log.error
		err = j.stageProgressHandler.SetJobStatus(sdk_v1.NewSetJobStatusReq(j.metadata.jobKey, sdk_v1.JobStatus_JobCompletionDoneWithErrors)) // TODO add an error
		return j
	}
	err = j.stageProgressHandler.SetJobStatus(sdk_v1.NewSetJobStatusReq(j.metadata.jobKey, sdk_v1.JobStatus_JobCompletionDone))
	return j
}

func (j *Job) Compensate(compensateDefinitionFn api.CompensateDefinitionFn) api.CompensateChain {
	if j.stageErr == nil {
		// TODO log.info can't execute the job's compensate stage because all stages ran successfully
		return j
	}

	err := j.stageProgressHandler.SetJobStatus(sdk_v1.NewSetJobStatusReq(j.metadata.jobKey, sdk_v1.JobStatus_JobCompensationStarted))
	if err != nil {
		// TODO log.error
		return j
	}

	compensationCtx := NewCompensationContext(j.clone())
	err = compensateDefinitionFn(compensationCtx)

	if err != nil {
		// TODO log.error
		err = j.stageProgressHandler.SetJobStatus(sdk_v1.NewSetJobStatusReq(j.metadata.jobKey, sdk_v1.JobStatus_JobCompensationDoneWithErrors)) // TODO add a reason fields to create a description for the error
		return j
	}
	err = j.stageProgressHandler.SetJobStatus(sdk_v1.NewSetJobStatusReq(j.metadata.jobKey, sdk_v1.JobStatus_JobCompensationDone))
	if err != nil {
		//TODO log.error
	}
	return j
}

func (j *Job) Canceled(cancelDefinitionFn api.CancelDefinitionFn) api.CanceledChain {
	err := j.stageProgressHandler.SetJobStatus(sdk_v1.NewSetJobStatusReq(j.metadata.jobKey, sdk_v1.JobStatus_JobCompensationStarted))
	if err != nil {
		// TODO log.error
		return j
	}

	cancellationCtx := NewCancellationContext(j.clone())
	err = cancelDefinitionFn(cancellationCtx)

	if err != nil {
		// TODO log.error
		err = j.stageProgressHandler.SetJobStatus(sdk_v1.NewSetJobStatusReq(j.metadata.jobKey, sdk_v1.JobStatus_JobCompensationDoneWithErrors)) // TODO add a reason fields to create a description for the error
		return j
	}
	err = j.stageProgressHandler.SetJobStatus(sdk_v1.NewSetJobStatusReq(j.metadata.jobKey, sdk_v1.JobStatus_JobCompensationDone))
	if err != nil {
		//TODO log.error
	}
	return j
}

func (j *Job) handleStageError(err error) api.StageChain {
	if se, ok := err.(api.StageError); ok {
		j.stageErr = se
	} else {
		j.stageErr = sdk_errors.NewStageError(err, sdk_errors.WithErrorType(sdk_v1.ErrorType_Failed))
	}
	// TODO log.error
	return j
}

func (j *Job) getStage(name string) (*sdk_v1.StageStatus, error) {
	return j.stageProgressHandler.Get(j.metadata.jobKey, name)
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

func (j *Job) updateStage(jobKey, name string, opts ...updateStageOption) error {
	req := &sdk_v1.SetStageStatusRequest{JobKey: jobKey, Name: name}
	for _, opt := range opts {
		req = opt(req)
	}
	return j.stageProgressHandler.Set(req)
}

func (j *Job) clone() *Job {
	clone := *j
	clone.stageErr = nil
	return &clone
}
