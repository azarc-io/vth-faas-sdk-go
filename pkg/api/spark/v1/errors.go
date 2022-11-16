package sdk_v1

import (
	"encoding/json"
	"errors"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	ErrStageDoesNotExist = errors.New("stage does not exists")
	ErrBindValueFailed   = errors.New("bind value failed")
	ErrVariableNotFound  = errors.New("variable not found")

	errorTypeToStageStatusMapper = map[ErrorType]StageStatus{
		ErrorType_ERROR_TYPE_RETRY:              StageStatus_STAGE_STATUS_FAILED,
		ErrorType_ERROR_TYPE_SKIP:               StageStatus_STAGE_STATUS_SKIPPED,
		ErrorType_ERROR_TYPE_CANCELLED:          StageStatus_STAGE_STATUS_CANCELLED,
		ErrorType_ERROR_TYPE_FAILED_UNSPECIFIED: StageStatus_STAGE_STATUS_FAILED,
	}
)

func ErrorTypeToStageStatusMapper(errType ErrorType) StageStatus {
	if err, ok := errorTypeToStageStatusMapper[errType]; ok {
		return err
	}
	return StageStatus_STAGE_STATUS_FAILED
}

type ErrorOption = func(err *Stage) *Stage

type Stage struct {
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

func NewStageError(err error, opts ...ErrorOption) *Stage {
	stg := &Stage{err: err}
	for _, opt := range opts {
		stg = opt(stg)
	}
	return stg
}

func (s *Stage) ErrorType() ErrorType {
	return s.errorType
}

func (s *Stage) Code() uint32 {
	return s.errorCode
}

func (s *Stage) Error() string {
	return s.err.Error()
}

func (s *Stage) Metadata() map[string]any {
	return s.metadata
}

func (s *Stage) ToErrorMessage() *Error {
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

func WithErrorType(errorType ErrorType) ErrorOption {
	return func(err *Stage) *Stage {
		err.errorType = errorType
		return err
	}
}

func WithErrorCode(code uint32) ErrorOption {
	return func(err *Stage) *Stage {
		err.errorCode = code
		return err
	}
}

func WithMetadata(metadata any) ErrorOption {
	return func(err *Stage) *Stage {
		err.parseMetadata(metadata)
		return err
	}
}

func WithRetry(times, backoffMillis uint) ErrorOption {
	return func(err *Stage) *Stage {
		err.retry = &RetryConfig{times, backoffMillis}
		err.errorType = ErrorType_ERROR_TYPE_RETRY
		return err
	}
}

func (s *Stage) parseMetadata(metadata any) {
	m := map[string]any{}
	if metadata != nil {
		mdBytes, _ := json.Marshal(metadata)
		_ = json.Unmarshal(mdBytes, &m)
	}
	s.metadata = m
}
