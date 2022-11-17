package spark_v1

import (
	"context"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/config"
)

/************************************************************************/
// TYPES
/************************************************************************/

type sparkWorker struct {
	config *config.Config
	chain  *chain
	ctx    context.Context
	opts   *sparkOpts
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
	metadata := NewSparkMetadata(context.Background(), jobKey, correlationID, transactionId, nil)
	return NewJobContext(metadata, w.opts)
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
		w.opts.variableHandler = NewIOHandler(grpcClient)
	}
	if w.opts.stageProgressHandler == nil {
		w.opts.stageProgressHandler = NewStageProgressHandler(grpcClient)
	}
	if w.opts.log == nil {
		w.opts.log = NewLogger()
	}
	return nil
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
