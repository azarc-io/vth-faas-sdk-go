package context

import (
	ctx "context"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	"github.com/rs/zerolog"
)

type Job struct {
	ctx                  ctx.Context
	metadata             JobMetadata
	stageProgressHandler api.StageProgressHandler
	variableHandler      api.VariableHandler
	log                  *zerolog.Logger
}

func NewJobContext(metadata api.Context, sph api.StageProgressHandler, vh api.VariableHandler, log *zerolog.Logger) api.JobContext {
	m := JobMetadata{ctx: metadata.Ctx(), jobKey: metadata.JobKey(), correlationId: metadata.CorrelationID(), transactionId: metadata.TransactionID(), payload: metadata.Payload(), lastActiveStage: metadata.LastActiveStage()}
	return &Job{metadata: m, stageProgressHandler: sph, variableHandler: vh, log: log}
}

func (j *Job) VariableHandler() api.VariableHandler {
	return j.variableHandler
}

func (j *Job) StageProgressHandler() api.StageProgressHandler {
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

func (j *Job) LastActiveStage() api.LastActiveStatus {
	return j.metadata.lastActiveStage
}

func (j *Job) SetVariables(stage string, variables ...*sdk_v1.Variable) error {
	return j.variableHandler.Set(j.metadata.jobKey, stage, variables...)
}

func (j *Job) GetVariables(stage string, names ...string) (*sdk_v1.Variables, error) {
	return j.variableHandler.Get(j.metadata.jobKey, stage, names...)
}

func (j *Job) Log() *zerolog.Logger {
	return j.log
}
