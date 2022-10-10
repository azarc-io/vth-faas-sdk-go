//go:generate mockgen -destination=../../internal/handlers/test/mock/mock_stageprogress.go -package=mock github.com/azarc-io/vth-faas-sdk-go/pkg/api StageProgressHandler
//go:generate mockgen -destination=../../internal/handlers/test/mock/mock_variable.go -package=mock github.com/azarc-io/vth-faas-sdk-go/pkg/api VariableHandler

package api

import (
	"context"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
)

type (
	Job interface {
		Initialize() error
		Execute(JobContext)
	}

	JobWorker interface {
		Run(ctx Context) StageError
	}

	VariableHandler interface {
		Get(jobKey, stage string, names ...string) ([]*sdk_v1.Variable, error)
		Set(jobKey, stage string, variables ...*sdk_v1.Variable) error
	}

	StageProgressHandler interface {
		Get(jobKey, name string) (*sdk_v1.StageStatus, error)
		Set(stageStatus *sdk_v1.SetStageStatusRequest) error
		GetResult(jobKey, name string) (*sdk_v1.StageResult, error)
		SetResult(resultResult *sdk_v1.SetStageResultRequest) error
		SetJobStatus(jobStatus *sdk_v1.SetJobStatusRequest) error
	}

	StageError interface {
		Error() string
		Code() uint32
		Metadata() map[string]any
		ErrorType() sdk_v1.ErrorType
		ToErrorMessage() *sdk_v1.Error
	}

	StageChain interface {
		Stage(name string, sdf StageDefinitionFn, options ...StageOption) StageChain
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
		Ctx() context.Context
		JobKey() string
		CorrelationID() string
		TransactionID() string
		Payload() any
	}

	JobContext interface {
		Stage(name string, sdf StageDefinitionFn, options ...StageOption) StageChain
		Err() StageError
	}

	StageContext interface {
		Context
		GetVariables(stage string, names ...string) ([]*sdk_v1.Variable, error)
	}

	CompletionContext interface {
		Context
		GetStageResult(stageName string) (*sdk_v1.StageResult, error)
		SetVariables(stage string, variables ...*sdk_v1.Variable) error
	}

	CompensationContext interface {
		Context
		Stage(name string, sdf StageDefinitionFn, options ...StageOption) StageChain
		GetVariables(stage string, names ...string) ([]*sdk_v1.Variable, error)
		SetVariables(stage string, variables ...*sdk_v1.Variable) error
	}

	CancelContext interface {
		Context
		Stage(name string, sdf StageDefinitionFn, options ...StageOption) StageChain
	}

	StageOptionParams interface {
		StageName() string
		StageProgressHandler() StageProgressHandler
		VariableHandler() VariableHandler
		Context() Context
	}

	StageDefinitionFn      = func(StageContext) (any, StageError)
	CompensateDefinitionFn = func(CompensationContext) StageError
	CancelDefinitionFn     = func(CancelContext) StageError
	CompletionDefinitionFn = func(CompletionContext) StageError
	StageOption            = func(StageOptionParams) StageError
)
