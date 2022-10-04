package context

import (
	"context"
	"errors"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
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

type Job struct {
	ctx                  context.Context
	metadata             JobMetadata
	stageProgressHandler api.StageProgressHandler
	variableHandler      api.VariableHandler
}

func NewJobContext(metadata JobMetadata, sph api.StageProgressHandler, vh api.VariableHandler) Job {
	return Job{metadata: metadata, stageProgressHandler: sph, variableHandler: vh}
}

func (j Job) JobKey() string {
	return j.metadata.jobKey
}

func (j Job) CorrelationID() string {
	return j.metadata.correlationId
}

func (j Job) TransactionID() string {
	return j.metadata.transactionId
}

func (j Job) Stage(name string, sdf api.StageDefinitionFn) api.StageChain {
	stage, err := j.stageProgressHandler.Get(name)
	if err != nil {
		// the request to start the job and run the stages was received,
		// but we can't call the server to retrieve the stage and execute the stage definition function
		// TODO panic(fmt.Sprintf("could not retrieve stage for job_key: %s; stage: %s; err: %s", j.metadata.jobKey, name, err.Error()))
		// or
		// TODO log.error
		return j
	}
	switch stage.Status {
	case sdk_v1.StageStatus_StagePending:
		stage.Status = sdk_v1.StageStatus_StageStarted
		err := j.stageProgressHandler.Set(stage)
		if err != nil {
			// TODO log.error
			return j
		}
		stageContext := NewStageContext(j)
		result, stageErr := sdf(stageContext)
		if stageErr != nil {
			stage = updateStageAttributesFromError(stage, stageErr)
			e := j.stageProgressHandler.Set(stage) // TODO configure retries via GRPC middleware => retry for ever https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/retry/examples_test.go
			if e != nil {
				// TODO log.error
			}
			return j
		}
		stage.Status = sdk_v1.StageStatus_StageCompleted
		err = j.stageProgressHandler.Set(stage)
		if err != nil {
			// TODO log.error
			return j
		}
		if result != nil {
			err := j.stageProgressHandler.SetResult(sdk_v1.NewStageResult(stage, result))
			if err != nil {
				// TODO log.error
				// panic(fmt.Sprintf("could not update stage result for job_key: %s; stage: %s; from status: %s; to status: %s; err: %s", stage.Job.Key, name, stage.Status, sdk_v1.StageStatus_StageCompleted, err.Error()))
			}
		}
		return j
	default:
		// StageFailed = 3;
		// StageSkipped = 4;
		// StageCanceled = 5;
		// TODO log.info stageName stageStatus
	}
	// TODO log.debug success
	return j
}

func (j Job) Complete(csd api.CompletionDefinitionFn) api.CompleteChain {
	job, err := j.stageProgressHandler.GetJob(j.JobKey())
	if err != nil {
		// TODO log.error
		return j
	}
	err = j.stageProgressHandler.SetJobStatus(job.Key, sdk_v1.JobStatus_JobCompletionRunning)
	if err != nil {
		// TODO log.error
		return j
	}
	completionCtx := NewCompleteContext(j)
	_, err = csd(completionCtx)
	if err != nil {
		// TODO log.error
		err = j.stageProgressHandler.SetJobStatus(job.Key, sdk_v1.JobStatus_JobCompletedWithErrors) // TODO add a reason fields to create a description for the error
		return j
	}
	err = j.stageProgressHandler.SetJobStatus(job.Key, sdk_v1.JobStatus_JobCompleted)
	return j
}

func (j Job) Compensate(cdf api.CompensateDefinitionFn) api.CompensateChain {
	job, err := j.stageProgressHandler.GetJob(j.JobKey())
	if err != nil {
		// TODO log.error
		return j
	}
	err = j.stageProgressHandler.SetJobStatus(job.Key, sdk_v1.JobStatus_JobCompensating)
	if err != nil {
		// TODO log.error
		return j
	}
	compensationCtx := NewCompensationContext(j)
	_, err = cdf(compensationCtx)
	if err != nil {
		// TODO log.error
		err = j.stageProgressHandler.SetJobStatus(job.Key, sdk_v1.JobStatus_JobCompensationCompletedWithErrors) // TODO add a reason fields to create a description for the error
		return j
	}
	err = j.stageProgressHandler.SetJobStatus(job.Key, sdk_v1.JobStatus_JobCompensationCompleted)
	return j
}

func (j Job) Canceled(cdf api.CancelDefinitionFn) api.CanceledChain {
	return j
}

func (j Job) Run() {
	//TODO do we really need this?
}

func updateStageAttributesFromError(stage *sdk_v1.Stage, err error) *sdk_v1.Stage {
	var stageError api.StageError
	if errors.As(err, &stageError) {
		if stageError.UpdateStatusTo() != nil {
			stage.Status = *stageError.UpdateStatusTo()
		}
		if stageError.Reason() != "" {
			stage.Reason = stageError.Reason()
		} else if err.Error() != "" {
			stage.Reason = stageError.Error()
		}
		stage.Retry = stageError.Retry()
		return stage
	}
	stage.Status = sdk_v1.StageStatus_StageFailed
	stage.Reason = err.Error()
	stage.Retry = true // TODO confirm if this is the right default value to that property
	return stage
}
