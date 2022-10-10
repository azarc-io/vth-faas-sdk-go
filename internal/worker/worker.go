package worker

import (
	"github.com/azarc-io/vth-faas-sdk-go/internal/context"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_errors "github.com/azarc-io/vth-faas-sdk-go/pkg/errors"
	"github.com/rs/zerolog"
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
	log                  zerolog.Logger
}

func NewJobWorker(job api.Job, options ...Option) (api.JobWorker, error) {
	jw := &JobWorker{job: job}
	for _, opt := range options {
		jw = opt(jw)
	}
	if err := jw.validate(); err != nil {
		return nil, err
	}
	return jw, nil
}

func (w JobWorker) validate() error {
	if w.variableHandler == nil {
		// TODO do not return an error use the default handler which is the GRPC handler
		return sdk_errors.VariableHandlerRequired
	}
	if w.stageProgressHandler == nil {
		// TODO do not return an error use the default handler which is the GRPC handler
		return sdk_errors.StageProgressHandlerRequired
	}
	return nil
}

func (w JobWorker) Run(metadata api.Context) api.StageError {
	jobContext := context.NewJobContext(metadata, w.stageProgressHandler, w.variableHandler)
	w.job.Execute(jobContext)
	return jobContext.Err()
}
