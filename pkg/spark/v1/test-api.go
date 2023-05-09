package sparkv1

import (
	"context"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
	"go.temporal.io/sdk/workflow"
)

/************************************************************************/
// IO
/************************************************************************/

type (
	IOState interface {
		GetVar(varName string) any
	}

	JobWorkflow interface {
		Run(ctx workflow.Context, jmd *JobMetadata) (*ExecuteSparkOutput, error)
		ExecuteStageActivity(ctx context.Context, req *ExecuteStageRequest) (Bindable, StageError)
		ExecuteCompleteActivity(ctx context.Context, req *ExecuteStageRequest) (*ExecuteSparkOutput, StageError)
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
		Inputs        ExecuteSparkInputs
		WorkflowId    string
		RunId         string
		TransactionId string
		CorrelationId string
		JobKey        string
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
