package spark_v1

import (
	ctx "context"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/internal/gen/azarc/sdk/spark/v1"
	"golang.org/x/net/context"
)

/************************************************************************/
// JOB CONTEXT
/************************************************************************/

type jobContext struct {
	metadata                *sparkMetadata
	stageProgressHandler    StageProgressHandler
	variableHandler         IOHandler
	log                     Logger
	delegateStageHandler    DelegateStageDefinitionFn
	delegateCompleteHandler DelegateCompleteDefinitionFn
}

func (j *jobContext) RequestMetadata() map[string]string {
	if j.metadata != nil {
		return j.metadata.perRequestMetadata
	}
	return map[string]string{}
}

func (j *jobContext) IOHandler() IOHandler {
	return j.variableHandler
}

func (j *jobContext) StageProgressHandler() StageProgressHandler {
	return j.stageProgressHandler
}

func (j *jobContext) Ctx() ctx.Context {
	return j.metadata.ctx
}

func (j *jobContext) JobKey() string {
	return j.metadata.jobKey
}

func (j *jobContext) CorrelationID() string {
	return j.metadata.correlationID
}

func (j *jobContext) TransactionID() string {
	return j.metadata.transactionID
}

func (j *jobContext) LastActiveStage() *sparkv1.LastActiveStage {
	return j.metadata.lastActiveStage
}

func (j *jobContext) Log() Logger {
	return j.log
}

func (j *jobContext) WithoutLastActiveStage() SparkContext {
	newCtx := *j
	md := *newCtx.metadata
	newCtx.metadata = &md
	newCtx.metadata.lastActiveStage = nil
	return &newCtx
}

func (j *jobContext) delegateStage() DelegateStageDefinitionFn {
	return j.delegateStageHandler
}

func (j *jobContext) delegateComplete() DelegateCompleteDefinitionFn {
	return j.delegateCompleteHandler
}

func NewJobContext(metadata Context, opts *sparkOpts) SparkContext {
	m := sparkMetadata{
		ctx:                metadata.Ctx(),
		jobKey:             metadata.JobKey(),
		correlationID:      metadata.CorrelationID(),
		transactionID:      metadata.TransactionID(),
		lastActiveStage:    metadata.LastActiveStage(),
		perRequestMetadata: metadata.RequestMetadata(),
	}
	return &jobContext{
		metadata:                &m,
		stageProgressHandler:    opts.stageProgressHandler,
		variableHandler:         opts.variableHandler,
		log:                     opts.log,
		delegateStageHandler:    opts.delegateStage,
		delegateCompleteHandler: opts.delegateComplete,
	}
}

/************************************************************************/
// METADATA
/************************************************************************/

type SparkMetadataOption func(sm *sparkMetadata)

type sparkMetadata struct {
	ctx                context.Context
	jobKey             string
	correlationID      string
	transactionID      string
	lastActiveStage    *sparkv1.LastActiveStage
	perRequestMetadata map[string]string
}

func (j *sparkMetadata) RequestMetadata() map[string]string {
	return j.perRequestMetadata
}

func (j *sparkMetadata) JobKey() string {
	return j.jobKey
}

func (j *sparkMetadata) CorrelationID() string {
	return j.correlationID
}

func (j *sparkMetadata) TransactionID() string {
	return j.transactionID
}

func (j *sparkMetadata) Ctx() context.Context {
	return j.ctx
}

func (j *sparkMetadata) LastActiveStage() *sparkv1.LastActiveStage {
	return j.lastActiveStage
}

func WithPerRequestMetadata(md map[string]string) SparkMetadataOption {
	return func(sm *sparkMetadata) {
		sm.perRequestMetadata = md
	}
}

func NewSparkMetadata(ctx context.Context, jobKey, correlationID, transactionID string, las *sparkv1.LastActiveStage, opts ...SparkMetadataOption) Context {
	sm := &sparkMetadata{
		ctx:             ctx,
		jobKey:          jobKey,
		correlationID:   correlationID,
		transactionID:   transactionID,
		lastActiveStage: las,
	}

	for _, opt := range opts {
		opt(sm)
	}

	return sm
}

func NewSparkMetadataFromGrpcRequest(ctx context.Context, req *sparkv1.ExecuteJobRequest) sparkMetadata {
	return sparkMetadata{
		ctx:             ctx,
		jobKey:          req.Key,
		correlationID:   req.CorrelationId,
		transactionID:   req.TransactionId,
		lastActiveStage: req.LastActiveStage,
	}
}

/************************************************************************/
// STAGE CONTEXT
/************************************************************************/

type stageContext struct {
	Context
	jobContext SparkContext
	name       string
}

func NewCompleteContext(ctx SparkContext, name string) CompleteContext {
	return stageContext{jobContext: ctx, Context: ctx, name: name}
}

func NewStageContext(ctx SparkContext, name string) StageContext {
	return stageContext{jobContext: ctx, Context: ctx, name: name}
}

func (sc stageContext) Inputs(names ...string) Inputs {
	return sc.jobContext.IOHandler().Inputs(sc.jobContext, names...)
}

func (sc stageContext) Input(name string) Input {
	return sc.jobContext.IOHandler().Input(sc.jobContext, name)
}

func (sc stageContext) StageResult(name string) Bindable {
	return sc.jobContext.StageProgressHandler().GetResult(sc.jobContext, name)
}

func (sc stageContext) Output(variables ...*Var) error {
	return sc.jobContext.IOHandler().Output(sc.jobContext, variables...)
}

func (sc stageContext) Name() string {
	return sc.name
}

func (sc stageContext) Log() Logger {
	return sc.jobContext.Log()
}

/************************************************************************/
// INIT CONTEXT
/************************************************************************/

type initContext struct {
	loader BindableConfig
	opts   *sparkOpts
}

func (i *initContext) Config() BindableConfig {
	if i.loader == nil {
		i.loader = newBindableConfig(i.opts)
	}
	return i.loader
}

func newInitContext(opts *sparkOpts) InitContext {
	return &initContext{opts: opts}
}
