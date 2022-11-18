package spark_v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/structpb"
	"time"
)

/************************************************************************/
// TYPES
/************************************************************************/

type ErrorOption = func(err *stageError) *stageError

type stageError struct {
	err       error
	errorType ErrorType
	errorCode uint32
	metadata  map[string]any
	retry     *RetryConfig
}

type RetryConfig struct {
	times         uint
	backoffMillis uint
}

/************************************************************************/
// ERRORS
/************************************************************************/

var (
	ErrStageDoesNotExist        = errors.New("stage does not exists")
	ErrBindValueFailed          = errors.New("bind value failed")
	ErrVariableNotFound         = errors.New("variable not found")
	ErrStageNotFoundInNodeChain = errors.New("stage not found in the node chain")
	ErrConditionalStageSkipped  = errors.New("conditional stage execution")
	ErrChainIsNotValid          = errors.New("chain is not valid")

	errorTypeToStageStatusMapper = map[ErrorType]StageStatus{
		ErrorType_ERROR_TYPE_RETRY:              StageStatus_STAGE_STATUS_FAILED,
		ErrorType_ERROR_TYPE_SKIP:               StageStatus_STAGE_STATUS_SKIPPED,
		ErrorType_ERROR_TYPE_CANCELLED:          StageStatus_STAGE_STATUS_CANCELLED,
		ErrorType_ERROR_TYPE_FAILED_UNSPECIFIED: StageStatus_STAGE_STATUS_FAILED,
	}
)

/************************************************************************/
// ERROR FACTORIES
/************************************************************************/

func newErrStageNotFoundInNodeChain(stage string) error {
	return fmt.Errorf("%w: %s", ErrStageNotFoundInNodeChain, stage)
}

func newErrConditionalStageSkipped(stageName string) error {
	return fmt.Errorf("%w: stage '%s' skipped", ErrConditionalStageSkipped, stageName)
}

func NewStageError(err error, opts ...ErrorOption) *stageError {
	stg := &stageError{err: err}
	for _, opt := range opts {
		stg = opt(stg)
	}
	return stg
}

/************************************************************************/
// STAGE ERROR ENVELOPE
/************************************************************************/

func (s *stageError) ErrorType() ErrorType {
	return s.errorType
}

func (s *stageError) Code() uint32 {
	return s.errorCode
}

func (s *stageError) Error() string {
	return s.err.Error()
}

func (s *stageError) Metadata() map[string]any {
	return s.metadata
}

func (s *stageError) ToErrorMessage() *Error {
	err := &Error{
		Error:     s.err.Error(),
		ErrorCode: s.errorCode,
		ErrorType: s.errorType,
	}
	if s.metadata != nil {
		err.Metadata, _ = structpb.NewValue(s.metadata)
	}
	if s.retry != nil {
		err.Retry = &RetryStrategy{Backoff: uint32(s.retry.backoffMillis), Count: uint32(s.retry.times)}
	}
	return err
}

/************************************************************************/
// STAGE ERROR OPTIONS
/************************************************************************/

func WithErrorType(errorType ErrorType) ErrorOption {
	return func(err *stageError) *stageError {
		err.errorType = errorType
		return err
	}
}

func WithErrorCode(code uint32) ErrorOption {
	return func(err *stageError) *stageError {
		err.errorCode = code
		return err
	}
}

func WithMetadata(metadata any) ErrorOption {
	return func(err *stageError) *stageError {
		err.parseMetadata(metadata)
		return err
	}
}

func WithRetry(times uint, backoffMillis time.Duration) ErrorOption {
	return func(err *stageError) *stageError {
		err.retry = &RetryConfig{times, uint(backoffMillis.Milliseconds())}
		err.errorType = ErrorType_ERROR_TYPE_RETRY
		return err
	}
}

func WithSkip() ErrorOption {
	return func(err *stageError) *stageError {
		err.errorType = ErrorType_ERROR_TYPE_SKIP
		return err
	}
}

func WithCancel() ErrorOption {
	return func(err *stageError) *stageError {
		err.errorType = ErrorType_ERROR_TYPE_CANCELLED
		err.metadata = map[string]any{"reason": "canceled in stage"}
		return err
	}
}

func WithFatal() ErrorOption {
	return func(err *stageError) *stageError {
		err.errorType = ErrorType_ERROR_TYPE_FATAL
		return err
	}
}

func (s *stageError) parseMetadata(metadata any) {
	m := map[string]any{}
	if metadata != nil {
		mdBytes, _ := json.Marshal(metadata)
		_ = json.Unmarshal(mdBytes, &m)
	}
	s.metadata = m
}

/************************************************************************/
// HELPERS
/************************************************************************/

func ErrorTypeToStageStatus(errType ErrorType) StageStatus {
	if err, ok := errorTypeToStageStatusMapper[errType]; ok {
		return err
	}
	return StageStatus_STAGE_STATUS_FAILED
}
