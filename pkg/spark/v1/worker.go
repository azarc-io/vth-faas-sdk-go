package sparkv1

import (
	"context"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

/************************************************************************/
// TYPES
/************************************************************************/

type sparkWorker struct {
	config    *Config
	chain     *SparkChain
	ctx       context.Context
	opts      *SparkOpts
	plugin    *sparkPlugin
	createdAt time.Time
	spark     Spark
	initOnce  sync.Once
	cancel    context.CancelFunc
}

/************************************************************************/
// Worker IMPLEMENTATION
/************************************************************************/

// Run runs the worker and waits for kill signals and then gracefully shuts down the worker
func (w *sparkWorker) Run() {
	// init the spark
	w.initIfRequired()

	// start plugin: expects plugin start to block
	w.startPlugin()

	// gracefully shutdown
	w.opts.log.Info("spark: start graceful shutdown")
	w.opts.log.Info("spark: stopping")
	w.spark.Stop()
	w.opts.log.Info("spark: stopped")

	w.opts.log.Info("plugin: stopping")
	w.plugin.stop()
	w.opts.log.Info("plugin: stopped")
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

	return nil
}

func (w *sparkWorker) loadConfiguration(opts *SparkOpts) {
	c, err := loadSparkConfig(opts)
	if err != nil {
		panic(err)
	}
	w.config = c
}

func (w *sparkWorker) startPlugin() {
	w.opts.log.Info("plugin: startup plugin: %+v", w.config)
	if err := w.plugin.start(); err != nil {
		panic(err)
	}
}

func (w *sparkWorker) initIfRequired() {
	w.initOnce.Do(func() {
		err := w.spark.Init(NewInitContext(w.opts))
		if err != nil {
			panic(err)
		}

		// register this runner as worker in temporal
		//tc := w.config.Config.Temporal
		log.Info().Msgf("config: %+v", w.config)
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
		opts:      &SparkOpts{},
		createdAt: time.Now(),
		spark:     spark,
	}
	for _, opt := range options {
		jw.opts = opt(jw.opts)
	}

	// load the configuration if available
	jw.loadConfiguration(jw.opts)

	// build the SparkChain
	builder := NewBuilder()
	spark.BuildChain(builder)
	chain := builder.BuildChain()

	// validate the SparkChain
	report := generateReportForChain(chain)

	jw.chain = chain

	if err := jw.validate(report); err != nil {
		return nil, err
	}

	jw.opts.log.Info("setting up plugin")
	jw.plugin = newSparkPlugin(ctx, jw.config, chain)

	return jw, nil
}
