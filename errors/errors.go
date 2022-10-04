package errors

import (
	"encoding/json"
	"errors"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
)

var (
	StageDoesNotExist            = errors.New("stage does not exists")
	BindValueFailed              = errors.New("bind value failed")
	VariableNotFound             = errors.New("variable not found")
	VariableHandlerRequired      = errors.New("a variable handler is required to create a new job worker")
	StageProgressHandlerRequired = errors.New("a stage progress handler is required to create a new job worker")
)

type Stage struct {
	err            error
	updateStatusTo *sdk_v1.StageStatus
	retry          bool
	errorCode      uint
	metadata       map[string]any
}

func NewStageError(err error, metadata any, reason string, status sdk_v1.StageStatus, retry bool) *Stage {
	stg := &Stage{err: err, reason: reason, updateStatusTo: &status, retry: retry}
	if err == nil && reason != "" {
		stg.err = errors.New(reason)
	}
	if reason == "" && err != nil {
		stg.reason = err.Error()
	}
	stg.parseMetadata(metadata)
	return stg
}

func (s *Stage) parseMetadata(metadata any) {
	m := map[string]any{}
	if metadata != nil {
		mdBytes, _ := json.Marshal(metadata) // TODO handle error - log at least
		_ = json.Unmarshal(mdBytes, &m)      // TODO handle error - log at least
	}
	s.metadata = m
}

func (s Stage) Error() string {
	if s.reason != "" {
		return s.reason
	}
	return s.err.Error()
}

func (s *Stage) Metadata() map[string]interface{} {
	return s.metadata
}

func (s *Stage) Reason() string {
	return s.reason
}

func (s *Stage) Retry() bool {
	return s.retry
}

func (s *Stage) UpdateStatusTo() *sdk_v1.StageStatus {
	return s.updateStatusTo
}

type StageErrors struct{}

var instance = &StageErrors{}

func New() *StageErrors {
	return instance
}

func (s StageErrors) Canceled(reason string) api.StageError {
	return NewStageError(nil, nil, reason, sdk_v1.StageStatus_StageCanceled, false)
}

func (s StageErrors) Fail(err error) api.StageError {
	return NewStageError(err, nil, "", sdk_v1.StageStatus_StageFailed, false)
}

func (s StageErrors) FailWithMetadata(err error, metadata any) api.StageError {
	return NewStageError(err, metadata, "", sdk_v1.StageStatus_StageFailed, false)
}

func (s StageErrors) FailWithRetry(err error) api.StageError {
	return NewStageError(err, nil, "", sdk_v1.StageStatus_StageFailed, true)
}

func (s StageErrors) FailWithReason(err error, reason string) api.StageError {
	return NewStageError(err, nil, reason, sdk_v1.StageStatus_StageFailed, false)
}

func (s StageErrors) FailWithReasonAndMetadata(err error, reason string, metadata any) api.StageError {
	return NewStageError(err, metadata, reason, sdk_v1.StageStatus_StageFailed, false)
}

type Option = func(err *Stage) *Stage

func WithRetry() Option {
	return func(err *Stage) *Stage {
		err.retry = true
		return err
	}
}

func WithStatus(status sdk_v1.StageStatus) Option {
	return func(err *Stage) *Stage {
		err.updateStatusTo = &status
		return err
	}
}

func WithReason(reason string) Option {
	return func(err *Stage) *Stage {
		err.reason = reason
		return err
	}
}

func WithReasonFromError(err error) Option {
	return func(e *Stage) *Stage {
		e.err = err
		e.reason = err.Error()
		return e
	}
}

func WithMetadata(metadata any) Option {
	return func(err *Stage) *Stage {
		err.parseMetadata(metadata)
		return err
	}
}

func NewFromOptions(opts ...Option) *Stage {
	stg := &Stage{}
	for _, opt := range opts {
		stg = opt(stg)
	}
	return stg
}
