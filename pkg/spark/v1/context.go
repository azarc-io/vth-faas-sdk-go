package spark_v1

import (
	ctx "context"
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

func (j *jobContext) LastActiveStage() *LastActiveStage {
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
		ctx:             metadata.Ctx(),
		jobKey:          metadata.JobKey(),
		correlationID:   metadata.CorrelationID(),
		transactionID:   metadata.TransactionID(),
		lastActiveStage: metadata.LastActiveStage(),
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

type sparkMetadata struct {
	ctx             context.Context
	jobKey          string
	correlationID   string
	transactionID   string
	lastActiveStage *LastActiveStage
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

func (j *sparkMetadata) LastActiveStage() *LastActiveStage {
	return j.lastActiveStage
}

func NewSparkMetadata(ctx context.Context, jobKey, correlationID, transactionID string, lastActiveStage *LastActiveStage) Context {
	return &sparkMetadata{
		ctx:             ctx,
		jobKey:          jobKey,
		correlationID:   correlationID,
		transactionID:   transactionID,
		lastActiveStage: lastActiveStage,
	}
}

func NewSparkMetadataFromGrpcRequest(ctx context.Context, req *ExecuteJobRequest) sparkMetadata {
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
	return sc.jobContext.IOHandler().Inputs(sc.JobKey(), names...)
}

func (sc stageContext) Input(name string) Input {
	return sc.jobContext.IOHandler().Input(sc.JobKey(), name)
}

func (sc stageContext) StageResult(name string) Bindable {
	return sc.jobContext.StageProgressHandler().GetResult(sc.JobKey(), name)
}

func (sc stageContext) Output(variables ...*Var) error {
	return sc.jobContext.IOHandler().Output(sc.JobKey(), variables...)
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
