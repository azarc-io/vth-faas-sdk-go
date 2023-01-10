package sparkv1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
	"time"
)

/************************************************************************/
// TYPES
/************************************************************************/

type ErrorOption = func(err *stageError) *stageError

type stageError struct {
	err      error
	metadata map[string]any
	retry    *RetryConfig
}

/************************************************************************/
// ERRORS
/************************************************************************/

var (
	ErrStageNotFoundInNodeChain = errors.New("stage not found in the Node SparkChain")
	ErrConditionalStageSkipped  = errors.New("conditional Stage execution")
	ErrChainIsNotValid          = errors.New("SparkChain is not valid")
)

var (
	MimeJsonError = codec.MimeTypeJson.WithType("error")
)

/************************************************************************/
// ERROR FACTORIES
/************************************************************************/

func newErrStageNotFoundInNodeChain(stage string) error {
	return fmt.Errorf("%w: %s", ErrStageNotFoundInNodeChain, stage)
}

func newErrConditionalStageSkipped(stageName string) error {
	return fmt.Errorf("%w: Stage '%s' skipped", ErrConditionalStageSkipped, stageName)
}

func NewStageError(err error, opts ...ErrorOption) StageError {
	stg := &stageError{err: err}
	for _, opt := range opts {
		stg = opt(stg)
	}
	return stg
}

/************************************************************************/
// STAGE ERROR ENVELOPE
/************************************************************************/
func (s *stageError) Error() string {
	return s.err.Error()
}

func (s *stageError) Metadata() map[string]any {
	return s.metadata
}

func (s *stageError) GetRetryConfig() *RetryConfig {
	return s.retry
}

/************************************************************************/
// STAGE ERROR OPTIONS
/************************************************************************/

func WithMetadata(metadata any) ErrorOption {
	return func(err *stageError) *stageError {
		err.parseMetadata(metadata)
		return err
	}
}

func WithRetry(times uint, backoffMultiplier uint, firstBackoffWait time.Duration) ErrorOption {
	//TODO Change to have retries
	return func(err *stageError) *stageError {
		err.retry = &RetryConfig{Times: times, BackoffMultiplier: backoffMultiplier, FirstBackoffWait: firstBackoffWait}
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
