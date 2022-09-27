package api

import (
	"time"
)

type (
	Job interface {
		Execute(JobContext)
	}

	JobContext interface {
		JobKey() string
		CorrelationID() string
		TransactionID() string
		Retries() any // encapsulate retries struct from temporal?
		Stage(name string, sdf StageDefinitionFn) StageChain
	}

	StageErrors interface {
		Canceled(reason string) StageError
		Fail(err error) StageError
		FailWithRetry(err error, delay time.Duration) StageError
		FailWithReason(err error, reason string, metadata map[string]interface{}) StageError
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
		Cancel(CancelDefinitionFn) CanceledChain
		Run()
	}

	CompensateChain interface {
		Cancel(CancelDefinitionFn) CanceledChain
		Complete(CompletionDefinitionFn) CompleteChain
		Run()
	}

	CanceledChain interface {
		Compensate(CompensateDefinitionFn) CompensateChain
		Complete(CompletionDefinitionFn) CompleteChain
		Run()
	}

	CompensationContext interface {
		Stage(name string, sdf StageDefinitionFn) StageChain
		WithStageStatus(names []string, value any) bool
		GetVariable(string) StageVariable
		SetVariable(name string, value any, mimeType string) error
	}

	CancelContext interface {
		Stage(name string, sdf StageDefinitionFn) StageChain
	}

	CompletionContext interface {
		GetStage(name string) Stage
		SetVariable(name string, value any, mimeType string) error
	}

	Stage interface {
		Raw() any
		BindValue(any) Stage // what we should return here?
		Context() CompletionContext
	}

	StageContext interface {
		GetVariable(string) StageVariable
	}

	StageVariable interface {
		Raw() any
		Bind(any) error
	}

	StageDefinitionFn      = func(StageContext) (any, StageError)
	CompensateDefinitionFn = func(CompensationContext) StageError
	CancelDefinitionFn     = func(CancelContext) StageError
	CompletionDefinitionFn = func(CompletionContext) StageError

	StageError interface {
		Error() string
	}
)
