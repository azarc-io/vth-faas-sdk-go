package spark_v1

import (
	"context"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/config"
)

/************************************************************************/
// TYPES
/************************************************************************/

type sparkWorker struct {
	config               *config.Config
	chain                *chain
	variableHandler      IOHandler
	stageProgressHandler StageProgressHandler
	log                  Logger
	ctx                  context.Context
}

/************************************************************************/
// Worker IMPLEMENTATION
/************************************************************************/

// Execute execute a single job
func (w *sparkWorker) Execute(metadata Context) StageError {
	jobContext := NewJobContext(metadata, w.stageProgressHandler, w.variableHandler, w.log)
	return w.chain.Execute(jobContext)
}

// LocalContext generates a context that can be used when calling Execute directly instead of through the api.
func (w *sparkWorker) LocalContext(jobKey, correlationID, transactionId string) Context {
	metadata := NewSparkMetadata(context.Background(), jobKey, correlationID, transactionId, nil)
	return NewJobContext(metadata, w.stageProgressHandler, w.variableHandler, w.log)
}

// Run runs the worker and waits for kill signals and then gracefully shuts down the worker
func (w *sparkWorker) Run() {
	//TODO implement me
	panic("implement me")
}

/************************************************************************/
// HELPERS
/************************************************************************/

func (w *sparkWorker) validate(report ChainReport) error {
	if len(report.Errors) > 0 {
		for _, err := range report.Errors {
			w.log.Error(err, "validation failed")
		}

		return ErrChainIsNotValid
	}

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

/************************************************************************/
// FACTORY
/************************************************************************/

func NewSparkWorker(ctx context.Context, spark Spark, options ...Option) (Worker, error) {
	jw := &sparkWorker{ctx: ctx}
	for _, opt := range options {
		jw = opt(jw)
	}

	// build the chain
	builder := NewBuilder()
	spark.BuildChain(builder)
	chain := builder.buildChain()

	// validate the chain
	report := GenerateReportForChain(chain)

	jw.chain = chain

	if err := jw.validate(report); err != nil {
		return nil, err
	}

	return jw, nil
}

/************************************************************************/
// OPTIONS
/************************************************************************/

type Option = func(je *sparkWorker) *sparkWorker

func WithIOHandler(vh IOHandler) Option {
	return func(jw *sparkWorker) *sparkWorker {
		jw.variableHandler = vh
		return jw
	}
}

func WithStageProgressHandler(sph StageProgressHandler) Option {
	return func(jw *sparkWorker) *sparkWorker {
		jw.stageProgressHandler = sph
		return jw
	}
}

func WithLog(log Logger) Option {
	return func(jw *sparkWorker) *sparkWorker {
		jw.log = log
		return jw
	}
}
