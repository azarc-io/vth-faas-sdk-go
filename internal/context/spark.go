package context

import (
	ctx "context"

	v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1"
)

type Job struct {
	ctx                  ctx.Context
	metadata             *SparkMetadata
	stageProgressHandler v1.StageProgressHandler
	variableHandler      v1.IOHandler
	log                  v1.Logger
}

func NewJobContext(metadata v1.Context, sph v1.StageProgressHandler, vh v1.IOHandler, log v1.Logger) v1.SparkContext {
	m := SparkMetadata{ctx: metadata.Ctx(), jobKey: metadata.JobKey(), correlationID: metadata.CorrelationID(), transactionID: metadata.TransactionID(), lastActiveStage: metadata.LastActiveStage()}
	return &Job{metadata: &m, stageProgressHandler: sph, variableHandler: vh, log: log}
}

func (j *Job) IOHandler() v1.IOHandler {
	return j.variableHandler
}

func (j *Job) StageProgressHandler() v1.StageProgressHandler {
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

func (j *Job) LastActiveStage() *v1.LastActiveStage {
	return j.metadata.lastActiveStage
}

func (j *Job) Log() v1.Logger {
	return j.log
}

func (j *Job) WithoutLastActiveStage() v1.SparkContext {
	newCtx := *j
	md := *newCtx.metadata
	newCtx.metadata = &md
	newCtx.metadata.lastActiveStage = nil
	return &newCtx
}
