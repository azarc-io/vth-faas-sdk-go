package worker

import (
	"github.com/azarc-io/vth-faas-sdk-go/internal/clients"
	"github.com/azarc-io/vth-faas-sdk-go/internal/context"
	grpc_handler "github.com/azarc-io/vth-faas-sdk-go/internal/handlers/grpc"
	"github.com/azarc-io/vth-faas-sdk-go/internal/job"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

func WithLog(log *zerolog.Logger) Option {
	return func(jw *JobWorker) *JobWorker {
		jw.log = log
		return jw
	}
}

type JobWorker struct {
	config               *config.Config
	chain                *job.Chain
	variableHandler      api.VariableHandler
	stageProgressHandler api.StageProgressHandler
	log                  *zerolog.Logger
}

func NewJobWorker(cfg *config.Config, chain *job.Chain, options ...Option) (api.JobWorker, error) {
	jw := &JobWorker{config: cfg, chain: chain}
	for _, opt := range options {
		jw = opt(jw)
	}
	if err := jw.validate(); err != nil {
		return nil, err
	}
	return jw, nil
}

func (w *JobWorker) validate() error {
	if w.variableHandler == nil {
		w.variableHandler = grpc_handler.NewVariableHandler()
	}
	if w.stageProgressHandler == nil {
		client, err := clients.CreateManagerServiceClient(w.config)
		if err != nil {
			return err
		}
		w.stageProgressHandler = grpc_handler.NewStageProgressHandler(client)
	}
	if w.log == nil {
		l := log.With().Str("module", "job_worker").CallerWithSkipFrameCount(3).Stack().Logger()
		w.log = &l
	}
	return nil
}

func (w *JobWorker) Run(metadata api.Context) api.StageError {
	jobContext := context.NewJobContext(metadata, w.stageProgressHandler, w.variableHandler)
	return w.chain.Execute(jobContext)
}
