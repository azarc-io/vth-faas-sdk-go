package worker

import (
	"github.com/azarc-io/vth-faas-sdk-go/internal/context"
	grpc_handler "github.com/azarc-io/vth-faas-sdk-go/internal/handlers/grpc"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/config"
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
	config               *config.Config
	job                  api.Job
	variableHandler      api.VariableHandler
	stageProgressHandler api.StageProgressHandler
	log                  zerolog.Logger
}

func NewJobWorker(cfg *config.Config, job api.Job, options ...Option) (api.JobWorker, error) {
	jw := &JobWorker{config: cfg, job: job}
	for _, opt := range options {
		jw = opt(jw)
	}
	if err := jw.validate(); err != nil {
		return nil, err
	}
	err := job.Initialize()
	if err != nil {
		return nil, err
	}
	return jw, nil
}

func (w JobWorker) validate() error {
	if w.variableHandler == nil {
		w.variableHandler = grpc_handler.NewGrpcVariableHandler()
	}
	if w.stageProgressHandler == nil {
		client, err := grpc_handler.CreateManagerServiceClient(w.config)
		if err != nil {
			return err
		}
		w.stageProgressHandler = grpc_handler.NewGrpcStageProgressHandler(client)
	}
	return nil
}

func (w JobWorker) Run(metadata api.Context) api.StageError {
	jobContext := context.NewJobContext(metadata, w.stageProgressHandler, w.variableHandler)
	w.job.Execute(jobContext)
	return jobContext.Err()
}
