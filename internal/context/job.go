package context

import (
	ctx "context"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
)

type Job struct {
	ctx                  ctx.Context
	metadata             JobMetadata
	stageProgressHandler sdk_v1.StageProgressHandler
	variableHandler      sdk_v1.VariableHandler
	log                  sdk_v1.Logger
}

func NewJobContext(metadata sdk_v1.Context, sph sdk_v1.StageProgressHandler, vh sdk_v1.VariableHandler, log sdk_v1.Logger) sdk_v1.SparkContext {
	m := JobMetadata{ctx: metadata.Ctx(), jobKey: metadata.JobKey(), correlationId: metadata.CorrelationID(), transactionId: metadata.TransactionID(), payload: metadata.Payload(), lastActiveStage: metadata.LastActiveStage()}
	return &Job{metadata: m, stageProgressHandler: sph, variableHandler: vh, log: log}
}

func (j *Job) VariableHandler() sdk_v1.VariableHandler {
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

func (j *Job) Payload() any {
	return j.metadata.payload
}

func (j *Job) LastActiveStage() sdk_v1.LastActiveStatus {
	return j.metadata.lastActiveStage
}

func (j *Job) SetVariables(stage string, variables ...*sdk_v1.Variable) error {
	return j.variableHandler.Set(j.metadata.jobKey, stage, variables...)
}

func (j *Job) GetVariables(stage string, names ...string) (*sdk_v1.Variables, error) {
	return j.variableHandler.Get(j.metadata.jobKey, stage, names...)
}

func (j *Job) Log() sdk_v1.Logger {
	return j.log
}
