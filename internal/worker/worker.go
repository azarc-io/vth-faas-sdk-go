package worker

import (
	"github.com/azarc-io/vth-faas-sdk-go/internal/clients"
	"github.com/azarc-io/vth-faas-sdk-go/internal/context"
	grpc_handler "github.com/azarc-io/vth-faas-sdk-go/internal/handlers/grpc"
	"github.com/azarc-io/vth-faas-sdk-go/internal/logger"
	"github.com/azarc-io/vth-faas-sdk-go/internal/spark"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/config"
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

func WithLog(log api.Logger) Option {
	return func(jw *JobWorker) *JobWorker {
		jw.log = log
		return jw
	}
}

type JobWorker struct {
	config               *config.Config
	chain                *spark.Chain
	variableHandler      api.VariableHandler
	stageProgressHandler api.StageProgressHandler
	log                  api.Logger
}

func NewSparkWorker(cfg *config.Config, chain *spark.Chain, options ...Option) (api.Worker, error) {
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
		w.log = logger.NewLogger()
	}
	return nil
}

func (w *JobWorker) Run(metadata api.Context) api.StageError {
	jobContext := context.NewJobContext(metadata, w.stageProgressHandler, w.variableHandler, w.log)
	return w.chain.Execute(jobContext)
}
