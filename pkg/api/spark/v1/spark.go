package spark_v1

import (
	ctx "context"
	"golang.org/x/net/context"
)

/************************************************************************/
// CONTEXT
/************************************************************************/

type job struct {
	ctx                     ctx.Context
	metadata                *sparkMetadata
	stageProgressHandler    StageProgressHandler
	variableHandler         IOHandler
	log                     Logger
	delegateStageHandler    DelegateStageDefinitionFn
	delegateCompleteHandler DelegateCompleteDefinitionFn
}

func NewJobContext(metadata Context, opts *sparkOpts) SparkContext {
	m := sparkMetadata{
		ctx:             metadata.Ctx(),
		jobKey:          metadata.JobKey(),
		correlationID:   metadata.CorrelationID(),
		transactionID:   metadata.TransactionID(),
		lastActiveStage: metadata.LastActiveStage(),
	}
	return &job{
		metadata:                &m,
		stageProgressHandler:    opts.stageProgressHandler,
		variableHandler:         opts.variableHandler,
		log:                     opts.log,
		delegateStageHandler:    opts.delegateStage,
		delegateCompleteHandler: opts.delegateComplete,
	}
}

func (j *job) IOHandler() IOHandler {
	return j.variableHandler
}

func (j *job) StageProgressHandler() StageProgressHandler {
	return j.stageProgressHandler
}

func (j *job) Ctx() ctx.Context {
	return j.ctx
}

func (j *job) JobKey() string {
	return j.metadata.jobKey
}

func (j *job) CorrelationID() string {
	return j.metadata.correlationID
}

func (j *job) TransactionID() string {
	return j.metadata.transactionID
}

func (j *job) LastActiveStage() *LastActiveStage {
	return j.metadata.lastActiveStage
}

func (j *job) Log() Logger {
	return j.log
}

func (j *job) WithoutLastActiveStage() SparkContext {
	newCtx := *j
	md := *newCtx.metadata
	newCtx.metadata = &md
	newCtx.metadata.lastActiveStage = nil
	return &newCtx
}

func (j *job) delegateStage() DelegateStageDefinitionFn {
	return j.delegateStageHandler
}

func (j *job) delegateComplete() DelegateCompleteDefinitionFn {
	return j.delegateCompleteHandler
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

func NewSparkMetadata(ctx context.Context, jobKey, correlationID, transactionID string, lastActiveStage *LastActiveStage) Context {
	return sparkMetadata{ctx: ctx, jobKey: jobKey, correlationID: correlationID, transactionID: transactionID, lastActiveStage: lastActiveStage}
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

func (j sparkMetadata) JobKey() string {
	return j.jobKey
}

func (j sparkMetadata) CorrelationID() string {
	return j.correlationID
}

func (j sparkMetadata) TransactionID() string {
	return j.transactionID
}

func (j sparkMetadata) Ctx() context.Context {
	return j.ctx
}

func (j sparkMetadata) LastActiveStage() *LastActiveStage {
	return j.lastActiveStage
}
