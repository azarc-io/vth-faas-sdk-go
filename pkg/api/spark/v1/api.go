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

	BuilderChain interface {
		ChainNode
	}

	ChainNode interface {
		ChainStage
	}

	ChainStage interface {
		Stage(name string, stageDefinitionFn StageDefinitionFn, options ...StageOption) ChainStage
		ChainCompensate
		ChainCancelled
		ChainComplete
	}

	ChainCancelledOrComplete interface {
		ChainCancelled
		ChainComplete
	}

	ChainCompensate interface {
		Compensate(newNode ChainNodeFinalizer) ChainCancelledOrComplete
	}

	ChainCancelled interface {
		Cancelled(newNode ChainNodeFinalizer) ChainComplete
	}

	ChainComplete interface {
		Complete(completeDefinitionFn CompleteDefinitionFn, options ...StageOption) ChainNodeFinalizer
	}

	ChainNodeFinalizer interface {
		build() *node
	}

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
		BuildChain(b Builder) ChainNodeFinalizer
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
)
