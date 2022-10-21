package v1

import (
	"github.com/azarc-io/vth-faas-sdk-go/internal/clients"
	"github.com/azarc-io/vth-faas-sdk-go/internal/context"
	grpc_handler "github.com/azarc-io/vth-faas-sdk-go/internal/handlers/grpc"
	"github.com/azarc-io/vth-faas-sdk-go/internal/logger"
	"github.com/azarc-io/vth-faas-sdk-go/internal/spark"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/config"
)

type Option = func(je *JobWorker) *JobWorker

func WithVariableHandler(vh sdk_v1.VariableHandler) Option {
	return func(jw *JobWorker) *JobWorker {
		jw.variableHandler = vh
		return jw
	}
}

func WithStageProgressHandler(sph sdk_v1.StageProgressHandler) Option {
	return func(jw *JobWorker) *JobWorker {
		jw.stageProgressHandler = sph
		return jw
	}
}

func WithLog(log sdk_v1.Logger) Option {
	return func(jw *JobWorker) *JobWorker {
		jw.log = log
		return jw
	}
}

type JobWorker struct {
	config               *config.Config
	chain                *spark.Chain
	variableHandler      sdk_v1.VariableHandler
	stageProgressHandler sdk_v1.StageProgressHandler
	log                  sdk_v1.Logger
}

func NewSparkWorker(cfg *config.Config, chain *spark.Chain, options ...Option) (sdk_v1.Worker, error) {
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

func (w *JobWorker) Execute(metadata sdk_v1.Context) sdk_v1.StageError {
	jobContext := context.NewJobContext(metadata, w.stageProgressHandler, w.variableHandler, w.log)
	return w.chain.Execute(jobContext)
}
