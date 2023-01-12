package sparkv1

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
	"github.com/pkg/errors"
)

/************************************************************************/
// TYPES
/************************************************************************/

type ErrorOption = func(err *stageError) *stageError

type stageError struct {
	stageName string
	err       error
	metadata  map[string]any
	retry     *RetryConfig
}

func (s *stageError) StackTrace() errors.StackTrace {
	if st, ok := s.err.(stackTracer); ok {
		return st.StackTrace()
	}

	// not stack tracable
	return nil
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
	if _, ok := err.(stackTracer); !ok {
		err = errors.WithStack(err)
	}

	stg := &stageError{err: err}
	for _, opt := range opts {
		stg = opt(stg)
	}
	return stg
}

/************************************************************************/
// STAGE ERROR ENVELOPE
/************************************************************************/

func (s *stageError) StageName() string {
	return s.stageName
}

func (s *stageError) Error() string {
	return s.err.Error()
}

func (s *stageError) Metadata() map[string]any {
	return s.metadata
}

func (s *stageError) GetRetryConfig() *RetryConfig {
	return s.retry
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

/*********************************************************************
Helpers
 **********************************************************************/

func getStackTrace(err stackTracer) []StackTraceItem {
	var stackTrace []StackTraceItem
	for _, frame := range err.StackTrace() {
		vals := strings.Split(fmt.Sprintf("%+v", frame), "\n")

		stackTrace = append(stackTrace, StackTraceItem{
			Type:     strings.TrimSpace(vals[0]),
			Filepath: strings.TrimSpace(vals[1]),
		})
	}
	return stackTrace
}
