package test

import (
	ctx "context"
	"github.com/azarc-io/vth-faas-sdk-go/internal/context"
	"github.com/azarc-io/vth-faas-sdk-go/internal/handlers/test/inmemory"
	"github.com/azarc-io/vth-faas-sdk-go/internal/spark"
	v1 "github.com/azarc-io/vth-faas-sdk-go/internal/worker/v1"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	"github.com/samber/lo"
	"testing"
)

func TestConditionalStageExecution(t *testing.T) {
	tests := []struct {
		name     string
		chainFn  func() (*spark.Chain, *stageBehaviour)
		assertFn func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler)
	}{
		{
			name: "should skip the second stage and only execute the first and third stages",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1", "stage2", "stage3")
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Stage("stage2", stageFn("stage2", sb), spark.WithStageStatus("stage1", sdk_v1.StageStatus_StageFailed)).
						Stage("stage3", stageFn("stage3", sb)).
						Build()).
					Build()
				if err != nil {
					t.Errorf("error creating spark node chain: %v", err)
				}
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				if sb.Executed("stage2") {
					t.Error("'stage2' should be skipped")
				}
				for _, stage := range []string{"stage1", "stage3"} {
					if !sb.Executed(stage) {
						t.Errorf("'%s' not executed", stage)
					}
				}
				stage2Status, err := sph.Get("jobKey", "stage2")
				if err != nil {
					t.Fatal(err)
				}
				if stage2Status != lo.ToPtr(sdk_v1.StageStatus_StageSkipped) {
					t.Errorf("'stage2' should be in 'skipped' status, got: %s", stage2Status)
				}
			},
		},
		{
			name: "should execute first stage and skip remaining 2 stages",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1", "stage2", "stage3")
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Stage("stage2", stageFn("stage2", sb), spark.WithStageStatus("stage1", sdk_v1.StageStatus_StageFailed)).
						Stage("stage3", stageFn("stage3", sb), spark.WithStageStatus("stage2", sdk_v1.StageStatus_StageCanceled)).
						Build()).
					Build()
				if err != nil {
					t.Errorf("error creating spark node chain: %v", err)
				}
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				if !sb.Executed("stage1") {
					t.Error("'stage1' should executed")
				}
				for _, stage := range []string{"stage2", "stage3"} {
					if !sb.Executed(stage) {
						t.Errorf("'%s' not executed", stage)
					}
					status, err := sph.Get("jobKey", stage)
					if err != nil {
						t.Fatal(err)
					}
					if status != lo.ToPtr(sdk_v1.StageStatus_StageSkipped) {
						t.Errorf("'stage2' should be in 'skipped' status, got: %s", status)
					}
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			chain, _ := test.chainFn() // sb
			stageProgressHandler := inmemory.NewStageProgressHandler(t)
			worker := v1.NewSparkTestWorker(t, chain, v1.WithVariableHandler(inmemory.NewVariableHandler(t, nil)), v1.WithStageProgressHandler(stageProgressHandler))
			err := worker.Execute(context.NewJobMetadata(ctx.Background(), "jobKey", "correlationId", "transactionId", nil))
			if err != nil {
				t.Error(err)
			}
		})
	}
}
