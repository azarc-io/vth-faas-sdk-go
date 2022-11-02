package context

import (
	ctx "context"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
)

type Job struct {
	ctx                  ctx.Context
	metadata             *JobMetadata
	stageProgressHandler sdk_v1.StageProgressHandler
	variableHandler      sdk_v1.IOHandler
	log                  sdk_v1.Logger
}

func NewJobContext(metadata sdk_v1.Context, sph sdk_v1.StageProgressHandler, vh sdk_v1.IOHandler, log sdk_v1.Logger) sdk_v1.SparkContext {
	m := JobMetadata{ctx: metadata.Ctx(), jobKey: metadata.JobKey(), correlationId: metadata.CorrelationID(), transactionId: metadata.TransactionID(), lastActiveStage: metadata.LastActiveStage()}
	return &Job{metadata: &m, stageProgressHandler: sph, variableHandler: vh, log: log}
}

func (j *Job) IOHandler() sdk_v1.IOHandler {
	return j.variableHandler
}

func (j *Job) StageProgressHandler() sdk_v1.StageProgressHandler {
	return j.stageProgressHandler
}

func (j *Job) Ctx() ctx.Context {
	return j.ctx
}

func (j *Job) JobKey() string {
	return j.metadata.jobKey
}

func (j *Job) CorrelationID() string {
	return j.metadata.correlationId
}

func (j *Job) TransactionID() string {
	return j.metadata.transactionId
}

func (j *Job) LastActiveStage() *sdk_v1.LastActiveStage {
	return j.metadata.lastActiveStage
}

func (j *Job) Log() sdk_v1.Logger {
	return j.log
}

func (j *Job) WithoutLastActiveStage() sdk_v1.SparkContext {
	newCtx := *j
	md := *newCtx.metadata
	newCtx.metadata = &md
	newCtx.metadata.lastActiveStage = nil
	return &newCtx
}
