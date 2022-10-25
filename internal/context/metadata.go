package context

import (
	"context"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
)

type JobMetadata struct {
	ctx             context.Context
	jobKey          string
	correlationId   string
	transactionId   string
	payload         any
	lastActiveStage *sdk_v1.LastActiveStage
}

func NewJobMetadata(ctx context.Context, jobKey, correlationId, transactionId string, lastActiveStage *sdk_v1.LastActiveStage) JobMetadata {
	return JobMetadata{ctx: ctx, jobKey: jobKey, correlationId: correlationId, transactionId: transactionId, lastActiveStage: lastActiveStage}
}

func NewJobMetadataFromGrpcRequest(ctx context.Context, req *sdk_v1.ExecuteJobRequest) JobMetadata {
	return JobMetadata{
		ctx:             ctx,
		jobKey:          req.Key,
		correlationId:   req.CorrelationId,
		transactionId:   req.TransactionId,
		lastActiveStage: req.LastActiveStage,
	}
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

func (j JobMetadata) LastActiveStage() *sdk_v1.LastActiveStage {
	return j.lastActiveStage
}
