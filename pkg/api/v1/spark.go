//go:generate mockgen -destination=../../internal/handlers/test/mock/mock_stageprogress.go -package=mock github.com/azarc-io/vth-faas-sdk-go/pkg/api StageProgressHandler
//go:generate mockgen -destination=../../internal/handlers/test/mock/mock_variable.go -package=mock github.com/azarc-io/vth-faas-sdk-go/pkg/api VariableHandler

package sdk_v1

import (
	"context"
)

type (
	Spark interface {
		Initialize() error
		Execute(SparkContext)
	}

	Worker interface {
		Execute(ctx Context) StageError
	}

	VariableHandler interface {
		Get(jobKey, stage string, names ...string) (*Variables, error)
		Set(jobKey, stage string, variables ...*Variable) error
	}

	StageProgressHandler interface {
		Get(jobKey, name string) (*StageStatus, error)
		Set(stageStatus *SetStageStatusRequest) error
		GetResult(jobKey, name string) (*StageResult, error)
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
		Payload() any
		LastActiveStage() LastActiveStatus
	}

	SparkContext interface {
		Context
		VariableHandler() VariableHandler
		StageProgressHandler() StageProgressHandler
		LastActiveStage() LastActiveStatus
		Log() Logger
	}

	LastActiveStatus interface {
		Name() string
		Status() StageStatus
	}

	StageContext interface {
		Context
		Inputs(stage string, names ...string) (*Variables, error)
		Input(stage string, names string) (*Variable, error)
		Log() Logger
	}

	CompleteContext interface {
		Context
		GetVariables(stage string, names ...string) (*Variables, error)
		Output(stage string, variables ...*Variable) error
		StageResult(name string) (*StageResult, error)
		Log() Logger
	}

	StageOptionParams interface {
		StageName() string
		StageProgressHandler() StageProgressHandler
		VariableHandler() VariableHandler
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
