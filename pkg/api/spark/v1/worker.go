package sdk_v1

import (
	"github.com/azarc-io/vth-faas-sdk-go/pkg/config"
)

type SparkWorker struct {
	config               *config.Config
	chain                *chain
	variableHandler      IOHandler
	stageProgressHandler StageProgressHandler
	log                  Logger
}

func (w *SparkWorker) Run() {
	//TODO implement me
	panic("implement me")
}

func (w *SparkWorker) Execute(metadata Context) StageError {
	jobContext := NewJobContext(metadata, w.stageProgressHandler, w.variableHandler, w.log)
	return w.chain.Execute(jobContext)
}

func (w *SparkWorker) validate() error {
	var grpcClient ManagerServiceClient
	if w.variableHandler == nil || w.stageProgressHandler == nil {
		var err error
		grpcClient, err = CreateManagerServiceClient(w.config)
		if err != nil {
			return err
		}
	}
	if w.variableHandler == nil {
		w.variableHandler = NewIOHandler(grpcClient)
	}
	if w.stageProgressHandler == nil {
		w.stageProgressHandler = NewStageProgressHandler(grpcClient)
	}
	if w.log == nil {
		w.log = NewLogger()
	}
	return nil
}

func NewSparkWorker(cfg *config.Config, chain *chain, options ...Option) (Worker, error) {
	jw := &SparkWorker{config: cfg, chain: chain}
	for _, opt := range options {
		jw = opt(jw)
	}
	if err := jw.validate(); err != nil {
		return nil, err
	}
	return jw, nil
}

type Option = func(je *SparkWorker) *SparkWorker

func WithIOHandler(vh IOHandler) Option {
	return func(jw *SparkWorker) *SparkWorker {
		jw.variableHandler = vh
		return jw
	}
}

func WithStageProgressHandler(sph StageProgressHandler) Option {
	return func(jw *SparkWorker) *SparkWorker {
		jw.stageProgressHandler = sph
		return jw
	}
}

func WithLog(log Logger) Option {
	return func(jw *SparkWorker) *SparkWorker {
		jw.log = log
		return jw
	}
}
