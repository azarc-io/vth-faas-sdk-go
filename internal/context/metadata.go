package context

import (
	"context"

	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
)

type SparkMetadata struct {
	ctx             context.Context
	jobKey          string
	correlationID   string
	transactionID   string
	lastActiveStage *sdk_v1.LastActiveStage
}

func NewSparkMetadata(ctx context.Context, jobKey, correlationID, transactionID string, lastActiveStage *sdk_v1.LastActiveStage) SparkMetadata {
	return SparkMetadata{ctx: ctx, jobKey: jobKey, correlationID: correlationID, transactionID: transactionID, lastActiveStage: lastActiveStage}
}

func NewSparkMetadataFromGrpcRequest(ctx context.Context, req *sdk_v1.ExecuteJobRequest) SparkMetadata {
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

func (j SparkMetadata) LastActiveStage() *sdk_v1.LastActiveStage {
	return j.lastActiveStage
}
