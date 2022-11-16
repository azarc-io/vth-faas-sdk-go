package sdk_v1

import (
	ctx "context"
	"golang.org/x/net/context"
)

/************************************************************************/
// CONTEXT
/************************************************************************/

type Job struct {
	ctx                  ctx.Context
	metadata             *SparkMetadata
	stageProgressHandler StageProgressHandler
	variableHandler      IOHandler
	log                  Logger
}

func NewJobContext(metadata Context, sph StageProgressHandler, vh IOHandler, log Logger) SparkContext {
	m := SparkMetadata{ctx: metadata.Ctx(), jobKey: metadata.JobKey(), correlationID: metadata.CorrelationID(), transactionID: metadata.TransactionID(), lastActiveStage: metadata.LastActiveStage()}
	return &Job{metadata: &m, stageProgressHandler: sph, variableHandler: vh, log: log}
}

func (j *Job) IOHandler() IOHandler {
	return j.variableHandler
}

func (j *Job) StageProgressHandler() StageProgressHandler {
	return j.stageProgressHandler
}

func (j *Job) Ctx() ctx.Context {
	return j.ctx
}

func (j *Job) JobKey() string {
	return j.metadata.jobKey
}

func (j *Job) CorrelationID() string {
	return j.metadata.correlationID
}

func (j *Job) TransactionID() string {
	return j.metadata.transactionID
}

func (j *Job) LastActiveStage() *LastActiveStage {
	return j.metadata.lastActiveStage
}

func (j *Job) Log() Logger {
	return j.log
}

func (j *Job) WithoutLastActiveStage() SparkContext {
	newCtx := *j
	md := *newCtx.metadata
	newCtx.metadata = &md
	newCtx.metadata.lastActiveStage = nil
	return &newCtx
}

/************************************************************************/
// METADATA
/************************************************************************/

type SparkMetadata struct {
	ctx             context.Context
	jobKey          string
	correlationID   string
	transactionID   string
	lastActiveStage *LastActiveStage
}

func NewSparkMetadata(ctx context.Context, jobKey, correlationID, transactionID string, lastActiveStage *LastActiveStage) SparkMetadata {
	return SparkMetadata{ctx: ctx, jobKey: jobKey, correlationID: correlationID, transactionID: transactionID, lastActiveStage: lastActiveStage}
}

func NewSparkMetadataFromGrpcRequest(ctx context.Context, req *ExecuteJobRequest) SparkMetadata {
	return SparkMetadata{
		ctx:             ctx,
		jobKey:          req.Key,
		correlationID:   req.CorrelationId,
		transactionID:   req.TransactionId,
		lastActiveStage: req.LastActiveStage,
	}
}

func (j SparkMetadata) JobKey() string {
	return j.jobKey
}

func (j SparkMetadata) CorrelationID() string {
	return j.correlationID
}

func (j SparkMetadata) TransactionID() string {
	return j.transactionID
}

func (j SparkMetadata) Ctx() context.Context {
	return j.ctx
}

func (j SparkMetadata) LastActiveStage() *LastActiveStage {
	return j.lastActiveStage
}
