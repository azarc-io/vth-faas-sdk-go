package sparkv1

import (
	"fmt"

	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
	"github.com/pkg/errors"
)

/************************************************************************/
// CONFIGURATION
/************************************************************************/

type ConfigType string

type ErrorCode string

const (
	ConfigTypeYaml ConfigType = "yaml"
	ConfigTypeJson ConfigType = "json"
)

const (
	ErrorCodeGeneric ErrorCode = "GENERIC"
)

/************************************************************************/
// ERRORS
/************************************************************************/

var (
	ErrTargetNotPointer            = errors.New("unable to set value of non-pointer")
	ErrUnableToBindUnknownMimeType = errors.New("unable to bind with unknown mime type")
)

/************************************************************************/
// BUILDER
/************************************************************************/

type (
	// Builder contract for the SparkChain builder
	Builder interface {
		NewChain(name string) BuilderChain
		ChainFinalizer
	}

	// BuilderChain the root of a SparkChain
	BuilderChain interface {
		ChainNode
	}

	// ChainNode a Node in the SparkChain
	ChainNode interface {
		ChainStage // must have at least 1 Stage
	}

	// ChainStage a Stage in the SparkChain Node
	ChainStage interface {
		Stage(name string, stageDefinitionFn StageDefinitionFn, options ...StageOption) ChainStageAny
	}

	// ChainStageAny allows defining more Stages and at least 1 of each Compensate, cancelled or Complete
	ChainStageAny interface {
		ChainStage
		ChainCompensate
		ChainCancelled
		ChainComplete
	}

	// ChainCancelledOrComplete allows defining only Cancel or completion
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

	// Chain finalizes a Node in the SparkChain, used internally to build a part of the SparkChain
	Chain interface {
		build() *Node
	}

	// ChainFinalizer finalizes the entire SparkChain, used internally to build the SparkChain
	ChainFinalizer interface {
		BuildChain() *SparkChain
	}
)

/************************************************************************/
// DATA APIS
/************************************************************************/

type (
	Gettable interface {
		Get(name string) Bindable
	}

	Bindable interface {
		Bind(a any) error
		GetValue() ([]byte, error)
		GetMimeType() string
	}

	BindableConfig interface {
		Bind(a any) error
		Raw() ([]byte, error)
	}

	Input interface {
		Bindable
	}

	Inputs interface {
		Get(name string) Bindable
	}

	bindable    Value
	BindableMap map[string]*bindable

	ExecuteSparkInputs BindableMap
	ExecuteSparkOutput struct {
		Outputs BindableMap        `json:"outputs,omitempty"`
		Error   *ExecuteSparkError `json:"error,omitempty"`
	}
	ExecuteSparkError struct {
		StageName    string           `json:"stage_name"`
		ErrorCode    ErrorCode        `json:"error_code"`
		ErrorMessage string           `json:"error_message,omitempty"`
		Metadata     map[string]any   `json:"metadata,omitempty"`
		StackTrace   []StackTraceItem `json:"stack_trace"`
	}

	SparkDataIO interface {
		GetStageResult(workflowId, runId, stageName string) (Bindable, error)
	}
)

func (b *bindable) Bind(a any) error {
	if b == nil || b.Value == nil {
		return nil
	}

	return errors.WithStack(codec.Decode(b.Value, a))
}

func (b *bindable) GetValue() ([]byte, error) {
	return b.Value, nil
}
func (b *bindable) GetMimeType() string {
	return b.MimeType
}

func NewBindable(value Value) *bindable {
	return &bindable{MimeType: value.MimeType, Value: value.Value}
}

func NewBindableValue(value any, mimeType string) *bindable {
	val, _ := codec.Encode(value)
	return &bindable{MimeType: mimeType, Value: val}
}

type errorBindable struct {
	err error
}

func (b *errorBindable) Bind(a any) error {
	return b.err
}

func (b *errorBindable) GetValue() ([]byte, error) {
	return nil, b.err
}
func (b *errorBindable) GetMimeType() string {
	return ""
}

func NewBindableError(err error) Bindable {
	return &errorBindable{err: err}
}

func (ese *ExecuteSparkError) Error() string {
	var stack []string
	for _, t := range ese.StackTrace {
		stack = append(stack, fmt.Sprintf("%s\n\t%s\n", t.Type, t.Filepath))
	}
	return fmt.Sprintf("%s\n%s", ese.ErrorMessage, stack)
}

/************************************************************************/
// CONTEXT
/************************************************************************/

type (
	Context interface {
		JobKey() string
		CorrelationID() string
		TransactionID() string
	}

	InitContext interface {
		Config() BindableConfig
	}

	StageContext interface {
		Context
		Input(names string) Input
		StageResult(name string) Bindable
		Log() Logger
		Name() string
	}

	CompleteContext interface {
		StageContext
		Output(variables ...*Var) error
		Name() string
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
		Init(ctx InitContext) error
		Stop()
	}
)

/************************************************************************/
// WORKER
/************************************************************************/

type (
	Worker interface {
		Run()
	}
)

/************************************************************************/
// ERRORS
/************************************************************************/

type (
	stackTracer interface {
		StackTrace() errors.StackTrace
	}

	StackTraceItem struct {
		Type     string `json:"type"`
		Filepath string `json:"filepath"`
	}

	StageError interface {
		stackTracer
		ErrorCode() ErrorCode
		StageName() string
		Error() string
		Metadata() map[string]any
		GetRetryConfig() *RetryConfig
	}
)

/************************************************************************/
// OPTIONS & PARAMS
/************************************************************************/

type (
	StageOptionParams interface {
		StageName() string
		Context() Context
	}

	StageDefinitionFn    = func(ctx StageContext) (any, StageError)
	CompleteDefinitionFn = func(ctx CompleteContext) StageError
	StageOption          = func(StageOptionParams) StageError
)
