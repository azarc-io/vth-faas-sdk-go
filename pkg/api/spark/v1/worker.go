package spark_v1

import (
	"context"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/signals"
)

/************************************************************************/
// TYPES
/************************************************************************/

type sparkWorker struct {
	config *Config
	chain  *chain
	ctx    context.Context
	opts   *sparkOpts
	server *Server
}

/************************************************************************/
// Worker IMPLEMENTATION
/************************************************************************/

// Execute execute a single job
func (w *sparkWorker) Execute(metadata Context) StageError {
	jobContext := NewJobContext(metadata, w.opts)
	return w.chain.execute(jobContext)
}

// LocalContext generates a context that can be used when calling Execute directly instead of through the api.
func (w *sparkWorker) LocalContext(jobKey, correlationID, transactionId string) Context {
	metadata := NewSparkMetadata(w.ctx, jobKey, correlationID, transactionId, nil)
	return NewJobContext(metadata, w.opts)
}

// Run runs the worker and waits for kill signals and then gracefully shuts down the worker
func (w *sparkWorker) Run() {
	// signal ch
	s := signals.SetupSignalHandler()
	// start server
	if w.server != nil {
		go func() {
			if err := w.server.Start(); err != nil {
				panic(err)
			}
		}()
		w.opts.log.Info("spark worker started, listening on: %s", w.config.ServerAddress())
	} else {
		w.opts.log.Info("spark worker started")
	}

	// wait for signals
	select {
	case <-s:
	case <-w.ctx.Done():
	}

	// gracefully shutdown
	w.opts.log.Info("gracefully shutting down spark")
	if w.server != nil {
		w.opts.log.Info("shutting down server")
		w.server.Stop()
		w.opts.log.Info("draining running sparks")
	}
}

/************************************************************************/
// HELPERS
/************************************************************************/

func (w *sparkWorker) validate(report ChainReport) error {
	if w.opts.log == nil {
		w.opts.log = NewLogger()
	}

	if len(report.Errors) > 0 {
		for _, err := range report.Errors {
			w.opts.log.Error(err, "validation failed")
		}

		return ErrChainIsNotValid
	}

	var grpcClient ManagerServiceClient
	if w.opts.variableHandler == nil || w.opts.stageProgressHandler == nil {
		var err error
		grpcClient, err = CreateManagerServiceClient(w.config)
		if err != nil {
			return err
		}
	}

	if w.opts.variableHandler == nil {
		w.opts.log.Info("setting up grpc i/o handler")
		w.opts.variableHandler = NewIOHandler(grpcClient)
	}

	if w.opts.stageProgressHandler == nil {
		w.opts.log.Info("setting up grpc progress handler")
		w.opts.stageProgressHandler = NewStageProgressHandler(grpcClient)
	}

	if w.config.Config.Server != nil && w.config.Config.Server.Enabled {
		w.opts.log.Info("setting up server")
		w.server = NewServer(w.config, w)
	}

	return nil
}

func (w *sparkWorker) loadConfiguration() {
	c, err := loadConfig()
	if err != nil {
		panic(err)
	}
	w.config = c
}

/************************************************************************/
// FACTORY
/************************************************************************/

func NewSparkWorker(ctx context.Context, spark Spark, options ...Option) (Worker, error) {
	jw := &sparkWorker{
		ctx:  ctx,
		opts: &sparkOpts{},
	}
	for _, opt := range options {
		jw.opts = opt(jw.opts)
	}

	// load the configuration if available
	jw.loadConfiguration()

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
