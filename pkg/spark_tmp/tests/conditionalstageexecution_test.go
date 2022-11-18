package tests

// TODO Move tests to /pkg/spark and get them working again
//import (
//	ctx "context"
//	"github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1/context"
//	"testing"
//
//	v12 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1"
//	v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/worker/v1"
//)
//
//func TestConditionalStageExecution(t *testing.T) {
//	tests := []struct {
//		name     string
//		chainFn  func() (*v12.BuilderChain, *stageBehaviour)
//		assertFn func(t *testing.T, sb *stageBehaviour, sph v12.StageProgressHandler)
//	}{
//		{
//			name: "should skip the second stage and only execute the first and third stages",
//			chainFn: func() (*v12.BuilderChain, *stageBehaviour) {
//				sb := NewStageBehaviour(t, "stage1", "stage2", "stage3")
//				chain, err := v12.NewChain(
//					v12.NewNode().
//						Stage("stage1", stageFn("stage1", sb)).
//						Stage("stage2", stageFn("stage2", sb), v12.WithStageStatus("stage1", v12.StageStatus_STAGE_STATUS_FAILED)).
//						Stage("stage3", stageFn("stage3", sb)).
//						Build()).
//					Build()
//				if err != nil {
//					t.Errorf("error creating spark node chain: %v", err)
//				}
//				return chain, sb
//			},
//			assertFn: func(t *testing.T, sb *stageBehaviour, sph v12.StageProgressHandler) {
//				if sb.Executed("stage2") {
//					t.Error("'stage2' should be skipped")
//				}
//				for _, stage := range []string{"stage1", "stage3"} {
//					if !sb.Executed(stage) {
//						t.Errorf("'%s' not executed", stage)
//					}
//				}
//				stage2Status, err := sph.Get("jobKey", "stage2")
//				if err != nil {
//					t.Fatal(err)
//				}
//				if *stage2Status != v12.StageStatus_STAGE_STATUS_SKIPPED {
//					t.Errorf("'stage2' should be in 'skipped' status, got: %s", stage2Status)
//				}
//			},
//		},
//		{
//			name: "should execute first stage and skip remaining 2 stages",
//			chainFn: func() (*v12.BuilderChain, *stageBehaviour) {
//				sb := NewStageBehaviour(t, "stage1", "stage2", "stage3")
//				chain, err := v12.NewChain(
//					v12.NewNode().
//						Stage("stage1", stageFn("stage1", sb)).
//						Stage("stage2", stageFn("stage2", sb), v12.WithStageStatus("stage1", v12.StageStatus_STAGE_STATUS_FAILED)).
//						Stage("stage3", stageFn("stage3", sb), v12.WithStageStatus("stage2", v12.StageStatus_STAGE_STATUS_CANCELLED)).
//						Build()).
//					Build()
//				if err != nil {
//					t.Errorf("error creating spark node chain: %v", err)
//				}
//				return chain, sb
//			},
//			assertFn: func(t *testing.T, sb *stageBehaviour, sph v12.StageProgressHandler) {
//				if !sb.Executed("stage1") {
//					t.Error("'stage1' should have executed")
//				}
//				for _, stage := range []string{"stage2", "stage3"} {
//					if sb.Executed(stage) {
//						t.Errorf("'%s' should not have executed", stage)
//					}
//					status, err := sph.Get("jobKey", stage)
//					if err != nil {
//						t.Fatal(err)
//					}
//					if *status != v12.StageStatus_STAGE_STATUS_SKIPPED {
//						t.Errorf("'%s' should be in 'skipped' status, got: %s", stage, status)
//					}
//				}
//			},
//		},
//	}
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			chain, sb := test.chainFn()
//			stageProgressHandler := v12.NewStageProgressHandler(t)
//			worker := v1.NewSparkTestWorker(t, chain, v12.WithIOHandler(v12.NewIOHandler(t)), v12.WithStageProgressHandler(stageProgressHandler))
//			err := worker.execute(context.NewSparkMetadata(ctx.Background(), "jobKey", "correlationId", "transactionId", nil))
//			if err != nil {
//				t.Error(err)
//			}
//			test.assertFn(t, sb, stageProgressHandler)
//		})
//	}
//}
