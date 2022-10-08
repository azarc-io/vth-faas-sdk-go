//go:generate mockgen -destination=../../internal/handlers/test/mock/mock_stageprogress.go -package=mock github.com/azarc-io/vth-faas-sdk-go/pkg/api StageProgressHandler
//go:generate mockgen -destination=../../internal/handlers/test/mock/mock_variable.go -package=mock github.com/azarc-io/vth-faas-sdk-go/pkg/api VariableHandler

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

	LocalWorker interface {
		Execute(ctx Context)
	}

	RemoteWorker interface {
		Start() error
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
		Get(jobKey, name string) (*sdk_v1.StageStatus, error)
		Set(stageStatus *sdk_v1.SetStageStatusRequest) error
		GetResult(jobKey, name string) (*sdk_v1.StageResult, error)
		SetResult(resultResult *sdk_v1.SetStageResultRequest) error
		SetJobStatus(jobStatus *sdk_v1.SetJobStatusRequest) error
	}

	JobContext interface {
		Stage(name string, sdf StageDefinitionFn) StageChain // add ...options to conditional stage execution
	}

	StageError interface {
		Error() string
		Code() uint32
		Metadata() map[string]any
		ErrorType() sdk_v1.ErrorType
		ToErrorMessage() *sdk_v1.Error
	}

	StageChain interface {
		Stage(name string, sdf StageDefinitionFn) StageChain
		Complete(CompletionDefinitionFn) CompleteChain
		Compensate(CompensateDefinitionFn) CompensateChain
		Canceled(CancelDefinitionFn) CanceledChain
	}

	CompleteChain interface {
		Compensate(CompensateDefinitionFn) CompensateChain
		Canceled(CancelDefinitionFn) CanceledChain
	}

	CompensateChain interface {
		Canceled(CancelDefinitionFn) CanceledChain
		Complete(CompletionDefinitionFn) CompleteChain
	}

	CanceledChain interface {
		Compensate(CompensateDefinitionFn) CompensateChain
		Complete(CompletionDefinitionFn) CompleteChain
	}

	Context interface {
		JobKey() string
		CorrelationID() string
		TransactionID() string
	}

	StageContext interface {
		Context
		GetVariable(string) (*sdk_v1.Variable, error)
	}

	CompletionContext interface {
		Context
		GetStage(jobKey, name string) (*sdk_v1.StageStatus, error)
		GetStageResult(jobKey, stageName string) (*sdk_v1.StageResult, error)
		SetVariable(variable *sdk_v1.Variable) error
	}

	CompensationContext interface {
		Context
		Stage(name string, sdf StageDefinitionFn) StageChain
		WithStageStatus(names []string, value any) bool
		GetVariable(string) (*sdk_v1.Variable, error)
		SetVariable(name string, value any, mimeType string) error
	}

	CancelContext interface {
		Context
		Stage(name string, sdf StageDefinitionFn) StageChain
	}

	StageDefinitionFn      = func(StageContext) (any, StageError)
	CompensateDefinitionFn = func(CompensationContext) StageError
	CancelDefinitionFn     = func(CancelContext) StageError
	CompletionDefinitionFn = func(CompletionContext) StageError
)
