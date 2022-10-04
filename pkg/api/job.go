package api

import (
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
)

type (
	Job interface {
		Execute(JobContext)
	}

	JobWorker interface {
		Run(payload any)
	}

	JobTestWorker interface {
		JobWorker
		Execute() //TODO
	}

	VariableHandler interface {
		Get(name string) (*sdk_v1.Variable, error)
		Set(variable *sdk_v1.Variable) error
	}

	StageProgressHandler interface {
		Get(name string) (*sdk_v1.Stage, error)
		Set(*sdk_v1.Stage) error
		GetResult(*sdk_v1.Stage) (*sdk_v1.StageResult, error)
		SetResult(result *sdk_v1.StageResult) error
		GetJob(jobKey string) (*sdk_v1.Job, error)
		SetJobStatus(jobKey string, status sdk_v1.JobStatus) error
	}

	JobContext interface {
		JobKey() string
		CorrelationID() string
		TransactionID() string
		Stage(name string, sdf StageDefinitionFn) StageChain
	}

	StageErrors interface {
		Canceled(reason string) StageError
		Fail(err error) StageError
		FailWithMetadata(err error, metadata any) StageError
		FailWithRetry(err error) StageError
		FailWithReason(err error, reason string) StageError
		FailWithReasonAndMetadata(err error, reason string, metadata any) StageError
	}

	StageChain interface {
		Stage(name string, sdf StageDefinitionFn) StageChain
		Complete(CompletionDefinitionFn) CompleteChain
		Compensate(CompensateDefinitionFn) CompensateChain
		Canceled(CancelDefinitionFn) CanceledChain
		Run()
	}

	CompleteChain interface {
		Compensate(CompensateDefinitionFn) CompensateChain
		Canceled(CancelDefinitionFn) CanceledChain
		Run()
	}

	CompensateChain interface {
		Canceled(CancelDefinitionFn) CanceledChain
		Complete(CompletionDefinitionFn) CompleteChain
		Run()
	}

	CanceledChain interface {
		Compensate(CompensateDefinitionFn) CompensateChain
		Complete(CompletionDefinitionFn) CompleteChain
		Run()
	}

	CompensationContext interface {
		Stage(name string, cdf CompensateDefinitionFn) StageChain
		WithStageStatus(names []string, value any) bool
		GetVariable(string) sdk_v1.Variable
		SetVariable(name string, value any, mimeType string) error
	}

	CancelContext interface {
		Stage(name string, cdf CancelDefinitionFn) StageChain
	}

	CompletionContext interface {
		GetStage(name string) (*sdk_v1.Stage, error)
		GetStageResult(stage *sdk_v1.Stage) (*sdk_v1.StageResult, error)
		SetVariable(variable *sdk_v1.Variable) error
	}

	StageContext interface {
		GetVariable(string) (*sdk_v1.Variable, error)
	}

	StageDefinitionFn      = func(StageContext) (any, StageError)
	CompensateDefinitionFn = func(CompensationContext) (any, StageError)
	CancelDefinitionFn     = func(CancelContext) (any, StageError)
	CompletionDefinitionFn = func(CompletionContext) (any, StageError)

	StageError interface {
		Error() string
		Metadata() map[string]any
		Reason() string
		Retry() bool
		UpdateStatusTo() *sdk_v1.StageStatus
	}
)
