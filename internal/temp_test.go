package internal

import (
	"github.com/azarc-io/vth-faas-sdk-go/internal/context"
	"github.com/azarc-io/vth-faas-sdk-go/internal/handlers"
	"github.com/azarc-io/vth-faas-sdk-go/internal/worker"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	"testing"
)

type Jobster struct {
}

func (j Jobster) Execute(ctx api.JobContext) {
	ctx.Stage("stage1", func(stageContext api.StageContext) (any, api.StageError) {
		return "something", nil
	}).Stage("stage2", func(stageContext api.StageContext) (any, api.StageError) {
		return "another result", nil
	}).Run()
}

func TestTmp(t *testing.T) {
	jobMetadata := context.NewJobMetadata("jobKey", "correlationId", "transactionId", []byte("payload"))

	jobWorker, err := worker.NewJobWorker(Jobster{},
		worker.WithStageProgressHandler(handlers.NewMockStageProgressHandler(&sdk_v1.Stage{
			Name:   "stage1",
			Status: sdk_v1.StageStatus_StagePending,
			JobId:  "jobId-1",
			Data:   nil,
			Reason: "",
		}, &sdk_v1.Stage{
			Name:   "stage2",
			Status: sdk_v1.StageStatus_Pending,
			JobId:  "jobId-2",
			Data:   nil,
			Reason: "",
		})),
		worker.WithVariableHandler(handlers.NewMockVariableHandler()))

	if err != nil {
		t.Fail()
	}

	jobWorker.Run(jobMetadata)
}
