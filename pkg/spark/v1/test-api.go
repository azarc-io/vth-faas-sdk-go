package sparkv1

import (
	"context"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
	"github.com/nats-io/nats.go/jetstream"
)

/************************************************************************/
// IO
/************************************************************************/

type (
	IOState interface {
		GetVar(varName string) any
	}

	JobWorkflow interface {
		Run(msg jetstream.Msg)
		ExecuteStageActivity(ctx context.Context, req *ExecuteStageRequest, io SparkDataIO) (Bindable, StageError)
		ExecuteCompleteActivity(ctx context.Context, req *ExecuteStageRequest, io SparkDataIO) (*ExecuteStageResponse, StageError)
	}

	StageTracker interface {
		GetStageResult(name string) (data any, mime codec.MimeType, err StageError)
		AssertStageCompleted(stageName string)
		AssertStageStarted(stageName string)
		AssertStageSkipped(stageName string)
		AssertStageCancelled(stageName string)
		AssertStageFailed(stageName string)
		AssertStageResult(stageName string, expectedStageResult any)
		AssertStageOrder(stageNames ...string)
	}

	ExecuteStageRequest struct {
		StageName     string
		TransactionId string
		CorrelationId string
		JobKey        string
		Inputs        map[string]Bindable
	}

	ExecuteStageResponse struct {
		Outputs BindableMap        `json:"outputs,omitempty"`
		Error   *ExecuteSparkError `json:"error,omitempty"`
	}

	Value struct {
		Value    []byte `json:"value"`
		MimeType string `json:"mime_type"`
	}
)

type StageStatus string

const (
	StageStatus_STAGE_PENDING   StageStatus = "STAGE_PENDING"
	StageStatus_STAGE_STARTED   StageStatus = "STAGE_STARTED"
	StageStatus_STAGE_COMPLETED StageStatus = "STAGE_COMPLETED"
	StageStatus_STAGE_FAILED    StageStatus = "STAGE_FAILED"
	StageStatus_STAGE_SKIPPED   StageStatus = "STAGE_SKIPPED"
	StageStatus_STAGE_CANCELED  StageStatus = "CANCELED"
)

type InternalStageTracker interface {
	SetStageResult(name string, value Bindable)
	SetStageStatus(name string, status StageStatus)
}
