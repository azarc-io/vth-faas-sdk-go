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
	lastActiveStage sdk_v1.LastActiveStatus
}

func NewJobMetadata(ctx context.Context, jobKey, correlationId, transactionId string, payload any) JobMetadata {
	return JobMetadata{ctx: ctx, jobKey: jobKey, correlationId: correlationId, transactionId: transactionId, payload: payload}
}

func NewJobMetadataFromGrpcRequest(ctx context.Context, req *sdk_v1.ExecuteJobRequest) JobMetadata {
	jm := JobMetadata{
		ctx:           ctx,
		jobKey:        req.Key,
		correlationId: req.CorrelationId,
		transactionId: req.TransactionId,
		payload:       nil, // TODO fix that
	}
	if req.LastActiveStage != nil {
		jm.lastActiveStage = LastActiveStatus{
			name:   req.LastActiveStage.Name,
			status: req.LastActiveStage.Status,
		}
	}
	return jm
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

func (j JobMetadata) LastActiveStage() sdk_v1.LastActiveStatus {
	return j.lastActiveStage
}

type LastActiveStatus struct {
	name   string
	status sdk_v1.StageStatus
}

func (l LastActiveStatus) Name() string {
	return l.name
}

func (l LastActiveStatus) Status() sdk_v1.StageStatus {
	return l.status
}
