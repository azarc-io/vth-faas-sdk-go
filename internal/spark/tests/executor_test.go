package tests

import (
	ctx "context"
	"errors"
	"github.com/azarc-io/vth-faas-sdk-go/internal/context"
	"github.com/azarc-io/vth-faas-sdk-go/internal/handlers/test/inmemory"
	"github.com/azarc-io/vth-faas-sdk-go/internal/spark"
	v1 "github.com/azarc-io/vth-faas-sdk-go/internal/worker/v1"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	sdk_errors "github.com/azarc-io/vth-faas-sdk-go/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSparkExecutor(t *testing.T) {
	tests := []struct {
		name            string
		chainFn         func() (*spark.Chain, *stageBehaviour)
		assertFn        func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler)
		errorType       *sdk_v1.ErrorType
		lastActiveStage *sdk_v1.LastActiveStage
		prepare         func(sph *inmemory.StageProgressHandler)
	}{
		{
			name: "should execute single stage spark",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1")
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.True(t, sb.Executed("stage1"), "'stage1' must have been executed")
				stage1Status, err := sph.Get("jobKey", "stage1")
				assert.Nil(t, err, "error retrieving spark status for stage1: %v", err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_StageCompleted), stage1Status, "'stage1' should be in 'completed' status, got: %s", stage1Status)
			},
		},
		{
			name: "when stage return an error the status must be failed",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1").Change("stage1", sdk_errors.NewStageError(errors.New("stage1")))
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.True(t, sb.Executed("stage1"), "'stage1' must have been executed")
				stage1Status, err := sph.Get("jobKey", "stage1")
				assert.Nil(t, err, "error retrieving spark status for stage1: %v", err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_StageFailed), stage1Status, "'stage1' should be in 'failed' status, got: %s", stage1Status)
			},
			errorType: lo.ToPtr(sdk_v1.ErrorType_Failed),
		},
		{
			name: "should execute complete stage",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1", "complete")
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Complete("complete", func(completeContext sdk_v1.CompleteContext) sdk_v1.StageError {
							sb.exec("complete")
							return sb.shouldErr("complete")
						}).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.True(t, sb.Executed("stage1"), "'stage1' must have been executed")
				stage1Status, err := sph.Get("jobKey", "stage1")
				assert.Nil(t, err, "error retrieving stage status for stage1: %v", err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_StageCompleted), stage1Status, "'stage1' should be in 'completed' status, got: %s", stage1Status)
				assert.True(t, sb.Executed("stage1"), "'stage1' must have been executed")
				completeStatus, err := sph.Get("jobKey", "complete")
				assert.Nil(t, err, "error retrieving stage status for complete: %v", err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_StageCompleted), completeStatus, "'complete' should be in 'completed' status, got: %s", completeStatus)
			},
		},
		{
			name: "should return an error if last active stage does not exist in the chain",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1")
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			lastActiveStage: &sdk_v1.LastActiveStage{Name: "non-existent"},
			errorType:       lo.ToPtr(sdk_v1.ErrorType_Failed),
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.False(t, sb.Executed("stage1"), "'stage1' must not have been executed")
			},
		},
		{
			name: "should skip stage1 return skip error and skip stage2 using conditional stage execution",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1", "stage2").
					Change("stage1",
						sdk_errors.NewStageError(errors.New("err-stage1"), sdk_errors.WithErrorType(sdk_v1.ErrorType_Skip)))
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Stage("stage2", stageFn("stage2", sb), spark.WithStageStatus("stage1", sdk_v1.StageStatus_StageCompleted)).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.True(t, sb.Executed("stage1"), "'stage1' must have been executed")
				stage1Status, err := sph.Get("jobKey", "stage1")
				assert.Nil(t, err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_StageSkipped), stage1Status, "'stage1' should be in 'skipped' status, got: %s", stage1Status)
				assert.False(t, sb.Executed("stage2"), "'stage2' must NOT have been executed")
				stage2Status, err := sph.Get("jobKey", "stage2")
				assert.Nil(t, err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_StageSkipped), stage2Status, "'stage1' should be in 'skipped' status, got: %s", stage2Status)
			},
		},
		{
			name: "should fail when a error occur during conditional stage execution evaluation",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1", "stage2")
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Stage("stage2", stageFn("stage2", sb), spark.WithStageStatus("stage1", sdk_v1.StageStatus_StageSkipped)).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.True(t, sb.Executed("stage1"), "'stage1' must have been executed")
				stage1Status, err := sph.Get("jobKey", "stage1")
				assert.Nil(t, err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_StageCompleted), stage1Status, "'stage1' should be in 'completed' status, got: %s", stage1Status)
				assert.False(t, sb.Executed("stage2"), "'stage2' must NOT have been executed")
			},
			errorType: lo.ToPtr(sdk_v1.ErrorType_Failed),
			prepare: func(sph *inmemory.StageProgressHandler) {
				sph.AddBehaviour().Set("stage2", errors.New("error updating status for stage 2"))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			chain, sb := test.chainFn()
			sph := inmemory.NewStageProgressHandler(t)
			if test.prepare != nil {
				test.prepare(sph)
			}
			zerolog.SetGlobalLevel(zerolog.Disabled)
			worker := v1.NewSparkTestWorker(t, chain, v1.WithIOHandler(inmemory.NewIOHandler(t)), v1.WithStageProgressHandler(sph))
			err := worker.Execute(context.NewJobMetadata(ctx.Background(), "jobKey", "correlationId", "transactionId", test.lastActiveStage))
			if err != nil && test.errorType == nil {
				t.Errorf("a unexpected error occured: %v", err)
			}
			if test.errorType != nil {
				if err == nil {
					t.Errorf("error '%s' is expected from chain execution, got none", test.errorType)
				} else if *test.errorType != err.ErrorType() {
					t.Errorf("error expected: %v; got: %v;", test.errorType, err.ErrorType())
				}
			}
			test.assertFn(t, sb, sph)
			sph.ResetBehaviour()
		})
	}
}
