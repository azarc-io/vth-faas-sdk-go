package spark_v1

import (
	"context"
)

//go:generate mockgen -destination=./test/mock_context.go -package spark_v1_mock github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1 Context
//go:generate mockgen -destination=./test/mock_stageprogress.go -package=spark_v1_mock github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1 StageProgressHandler

/************************************************************************/
// BUILDER
/************************************************************************/

type (
	// Builder contract for the chain builder
	Builder interface {
		NewChain(name string) BuilderChain
		ChainFinalizer
	}

	// BuilderChain the root of a chain
	BuilderChain interface {
		ChainNode
	}

	// ChainNode a node in the chain
	ChainNode interface {
		ChainStage // must have at least 1 stage
	}

	// ChainStage a stage in the chain node
	ChainStage interface {
		Stage(name string, stageDefinitionFn StageDefinitionFn, options ...StageOption) ChainStageAny
	}

	// ChainStageAny allows defining more stages and at least 1 of each compensate, cancelled or complete
	ChainStageAny interface {
		ChainStage
		ChainCompensate
		ChainCancelled
		ChainComplete
	}

	// ChainCancelledOrComplete allows defining only cancel or completion
	ChainCancelledOrComplete interface {
		ChainCancelled
		ChainComplete
	}

	// ChainCompensate contract the builder must implement for compensation
	ChainCompensate interface {
		Compensate(newNode Chain) ChainCancelledOrComplete
	}

	// ChainCancelled contract the builder must implement for cancellation
	ChainCancelled interface {
		Cancelled(newNode Chain) ChainComplete
	}

	// ChainComplete contract the builder must implement for completion
	ChainComplete interface {
		Complete(completeDefinitionFn CompleteDefinitionFn, options ...StageOption) Chain
	}

	// Chain finalizes a node in the chain, used internally to build a part of the chain
	Chain interface {
		build() *node
	}

	// ChainFinalizer finalizes the entire chain, used internally to build the chain
	ChainFinalizer interface {
		buildChain() *chain
	}
)

/************************************************************************/
// IO
/************************************************************************/

type (
	IOHandler interface {
		Inputs(jobKey string, names ...string) *Inputs
		Input(jobKey, name string) *Input
		Output(jobKey string, variables ...*Var) error
	}
)

/************************************************************************/
// PROGRESS
/************************************************************************/

type (
	StageProgressHandler interface {
		Get(jobKey, name string) (*StageStatus, error)
		Set(stageStatus *SetStageStatusRequest) error
		GetResult(jobKey, name string) *Result
		SetResult(resultResult *SetStageResultRequest) error
		SetJobStatus(jobStatus *SetJobStatusRequest) error
	}
)

/************************************************************************/
// CONTEXT
/************************************************************************/

type (
	Context interface {
		Ctx() context.Context
		JobKey() string
		CorrelationID() string
		TransactionID() string
		LastActiveStage() *LastActiveStage
	}

	SparkContext interface {
		Context
		IOHandler() IOHandler
		StageProgressHandler() StageProgressHandler
		LastActiveStage() *LastActiveStage
		Log() Logger
		WithoutLastActiveStage() SparkContext
		delegateStage() DelegateStageDefinitionFn
		delegateComplete() DelegateCompleteDefinitionFn
	}

	StageContext interface {
		Context
		Inputs(names ...string) *Inputs
		Input(names string) *Input
		StageResult(name string) *Result
		Log() Logger
	}

	CompleteContext interface {
		StageContext
		Output(variables ...*Var) error
	}
)

/************************************************************************/
// LOGGING
/************************************************************************/

type (
	Logger interface {
		Info(format string, v ...any)
		Warn(format string, v ...any)
		Debug(format string, v ...any)
		Error(err error, format string, v ...any)
		AddFields(k string, v any) Logger
	}
)

/************************************************************************/
// SPARK
/************************************************************************/

type (
	// Spark the contract a developer must implement in order to be accepted by a worker
	Spark interface {
		BuildChain(b Builder) Chain
	}
)

/************************************************************************/
// WORKER
/************************************************************************/

type (
	Worker interface {
		Execute(ctx Context) StageError
		Run()
		LocalContext(jobKey, correlationID, transactionId string) Context
	}
)

/************************************************************************/
// ERRORS
/************************************************************************/

type (
	StageError interface {
		Error() string
		Code() uint32
		Metadata() map[string]any
		ErrorType() ErrorType
		ToErrorMessage() *Error
	}
)

/************************************************************************/
// OPTIONS & PARAMS
/************************************************************************/

type (
	StageOptionParams interface {
		StageName() string
		StageProgressHandler() StageProgressHandler
		IOHandler() IOHandler
		Context() Context
	}

	StageDefinitionFn    = func(ctx StageContext) (any, StageError)
	CompleteDefinitionFn = func(ctx CompleteContext) StageError
	StageOption          = func(StageOptionParams) StageError

	DelegateStageDefinitionFn    = func(ctx StageContext, cb StageDefinitionFn) (any, StageError)
	DelegateCompleteDefinitionFn = func(ctx CompleteContext, cb CompleteDefinitionFn) StageError
)
