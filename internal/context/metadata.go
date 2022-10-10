package context

import (
	"context"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
)

type JobMetadata struct {
	ctx           context.Context
	jobKey        string
	correlationId string
	transactionId string
	payload       any
}

func NewJobMetadata(ctx context.Context, jobKey, correlationId, transactionId string, payload any) JobMetadata {
	return JobMetadata{ctx, jobKey, correlationId, transactionId, payload}
}

func (j JobMetadata) JobKey() string {
	return j.jobKey
}

func (j JobMetadata) CorrelationID() string {
	return j.correlationId
}

func (j JobMetadata) TransactionID() string {
	return j.transactionId
}

func (j JobMetadata) Payload() any {
	return j.payload
}

func (j JobMetadata) Ctx() context.Context {
	return j.ctx
}

type Job struct {
	ctx                  context.Context
	metadata             JobMetadata
	stageProgressHandler api.StageProgressHandler
	variableHandler      api.VariableHandler
	stageErr             api.StageError
}
