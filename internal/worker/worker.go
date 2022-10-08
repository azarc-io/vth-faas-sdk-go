package worker

import (
	"github.com/azarc-io/vth-faas-sdk-go/internal/context"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_errors "github.com/azarc-io/vth-faas-sdk-go/pkg/errors"
	"testing"
)

type Option = func(je *JobWorker) *JobWorker

func WithVariableHandler(vh api.VariableHandler) Option {
	return func(jw *JobWorker) *JobWorker {
		jw.variableHandler = vh
		return jw
	}
}

func WithStageProgressHandler(sph api.StageProgressHandler) Option {
	return func(jw *JobWorker) *JobWorker {
		jw.stageProgressHandler = sph
		return jw
	}
}

type JobWorker struct {
	job                  api.Job
	variableHandler      api.VariableHandler
	stageProgressHandler api.StageProgressHandler
}

func NewJobWorker(job api.Job, options ...Option) (*JobWorker, error) {
	jw := &JobWorker{job: job}
	for _, opt := range options {
		jw = opt(jw)
	}
	if err := jw.validate(); err != nil {
		return nil, err
	}
	return jw, nil
}

// TODO return a test worker and expose the Execute methd
func NewTestJobWorker(t *testing.T, job api.Job, options ...Option) (*JobTestWorker, error) {
	return nil, nil
}

func (j JobWorker) validate() error {
	if j.variableHandler == nil {
		// TODO do not return an error use the default handler which is the GRPC handler
		return sdk_errors.VariableHandlerRequired
	}
	if j.stageProgressHandler == nil {
		// TODO do not return an error use the default handler which is the GRPC handler
		return sdk_errors.StageProgressHandlerRequired
	}
	return nil
}

// TODO need to hide that from the user in production
// can only be called by the user in testing
func (j JobWorker) Run(metadata api.Context) {
	jobContext := context.NewJobContext(metadata, j.stageProgressHandler, j.variableHandler)
	j.job.Execute(jobContext)
}
