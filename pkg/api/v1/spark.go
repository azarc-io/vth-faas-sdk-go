//go:generate mockgen -destination=../../internal/handlers/test/mock/mock_stageprogress.go -package=mock github.com/azarc-io/vth-faas-sdk-go/pkg/api StageProgressHandler
//go:generate mockgen -destination=../../internal/handlers/test/mock/mock_variable.go -package=mock github.com/azarc-io/vth-faas-sdk-go/pkg/api VariableHandler

package sdk_v1

import (
	"context"

	"github.com/azarc-io/vth-faas-sdk-go/internal/handlers"
)

type (
	Spark interface {
		Initialize() error
		Execute(SparkContext)
	}

	Worker interface {
		Execute(ctx Context) StageError
	}

	IOHandler interface {
		Inputs(jobKey string, names ...string) *Inputs
		Input(jobKey, name string) *Input
		Output(jobKey string, variables ...*handlers.Variable) error
	}

	StageProgressHandler interface {
		Get(jobKey, name string) (*StageStatus, error)
		Set(stageStatus *SetStageStatusRequest) error
		GetResult(jobKey, name string) *Result
		SetResult(resultResult *SetStageResultRequest) error
		SetJobStatus(jobStatus *SetJobStatusRequest) error
	}

	StageError interface {
		Error() string
		Code() uint32
		Metadata() map[string]any
		ErrorType() ErrorType
		ToErrorMessage() *Error
	}

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
		Output(variables ...*handlers.Variable) error
	}

	StageOptionParams interface {
		StageName() string
		StageProgressHandler() StageProgressHandler
		IOHandler() IOHandler
		Context() Context
	}

	StageDefinitionFn    = func(StageContext) (any, StageError)
	CompleteDefinitionFn = func(CompleteContext) StageError
	StageOption          = func(StageOptionParams) StageError

	Logger interface {
		Info(format string, v ...any)
		Warn(format string, v ...any)
		Debug(format string, v ...any)
		Error(err error, format string, v ...any)
		AddFields(k string, v any) Logger
	}
)
