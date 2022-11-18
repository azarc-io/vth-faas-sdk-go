package spark_v1

import (
	"context"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/healthz"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/signals"
	"net/http"
	"sync"
	"time"
)

/************************************************************************/
// TYPES
/************************************************************************/

type sparkWorker struct {
	config    *config
	chain     *chain
	ctx       context.Context
	opts      *sparkOpts
	server    *server
	createdAt time.Time
	health    *healthz.Checker
	spark     Spark
	initOnce  sync.Once
	cancel    context.CancelFunc
}

/************************************************************************/
// Worker IMPLEMENTATION
/************************************************************************/

// Execute execute a single job
// TODO this should not be exposed once all Azarc projects have been consolidated into Verathread
func (w *sparkWorker) Execute(metadata Context) StageError {
	// init the spark
	w.initIfRequired()

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
	w.startServer()

	// init the spark
	w.initIfRequired()

	// wait for signals
	select {
	case <-s:
		w.cancel()
	case <-w.ctx.Done():
	}

	// gracefully shutdown
	w.opts.log.Info("gracefully shutting down spark")
	if w.server != nil {
		w.opts.log.Info("shutting down server")
		w.server.stop()
	}

	w.spark.Stop()
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
		w.opts.variableHandler = newGrpcIOHandler(grpcClient)
	}

	if w.opts.stageProgressHandler == nil {
		w.opts.log.Info("setting up grpc progress handler")
		w.opts.stageProgressHandler = newGrpcStageProgressHandler(grpcClient)
	}

	if w.config.Config.Server != nil && w.config.Config.Server.Enabled {
		w.opts.log.Info("setting up server")
		w.server = newServer(w.config, w)
	}

	// TODO support TLS once support for platforms other than kubernetes are added to Verathread
	if w.config.Config.Health != nil && w.config.Config.Health.Enabled {
		w.opts.log.Info("setting up healthz")
		w.health = healthz.NewChecker(&healthz.Config{
			RuntimeTTL: time.Second * 5,
		})

		go func() {
			http.Handle("/healthz", w.health.Handler())

			// nosemgrep
			if err := http.ListenAndServe(w.config.healthBindTo(), nil); err != nil { // nosemgrep
				panic(err)
			}
		}()
	}

	return nil
}

func (w *sparkWorker) loadConfiguration() {
	c, err := loadSparkConfig()
	if err != nil {
		panic(err)
	}
	w.config = c
}

func (w *sparkWorker) startServer() {
	if w.server != nil {
		go func() {
			if err := w.server.start(); err != nil {
				panic(err)
			}
		}()
		w.opts.log.Info("spark worker started in %v, listening on: %s", time.Since(w.createdAt), w.config.serverAddress())
	} else {
		w.opts.log.Info("spark worker started in %v", time.Since(w.createdAt))
	}
}

func (w *sparkWorker) initIfRequired() {
	w.initOnce.Do(func() {
		err := w.spark.Init(newInitContext())
		if err != nil {
			panic(err)
		}
	})
}

/************************************************************************/
// FACTORY
/************************************************************************/

func NewSparkWorker(ctx context.Context, spark Spark, options ...Option) (Worker, error) {
	wrappedCtx, cancel := context.WithCancel(ctx)
	var jw = &sparkWorker{
		ctx:       wrappedCtx,
		cancel:    cancel,
		opts:      &sparkOpts{},
		createdAt: time.Now(),
		spark:     spark,
	}
	for _, opt := range options {
		jw.opts = opt(jw.opts)
	}

	// load the configuration if available
	jw.loadConfiguration()

	// build the chain
	builder := newBuilder()
	spark.BuildChain(builder)
	chain := builder.buildChain()

	// validate the chain
	report := generateReportForChain(chain)

	jw.chain = chain

	if err := jw.validate(report); err != nil {
		return nil, err
	}

	return jw, nil
}
