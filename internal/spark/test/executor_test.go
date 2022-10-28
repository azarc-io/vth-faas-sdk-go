package test

import (
	ctx "context"
	"github.com/azarc-io/vth-faas-sdk-go/internal/context"
	"github.com/azarc-io/vth-faas-sdk-go/internal/handlers/test/inmemory"
	"github.com/azarc-io/vth-faas-sdk-go/internal/spark"
	v1 "github.com/azarc-io/vth-faas-sdk-go/internal/worker/v1"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSparkExecutor(t *testing.T) {
	tests := []struct {
		name     string
		chainFn  func() (*spark.Chain, *stageBehaviour)
		assertFn func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler)
	}{
		{
			name: "should execute without errors a single stage spark",
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
				assert.True(t, sb.Executed("stage1"), "")
				stage1Status, err := sph.Get("jobKey", "stage1")
				assert.Nil(t, err, "error retrieving spark status for stage1: %v", err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_StageCompleted), "'stage1' should be in 'completed' status, got: %s", stage1Status)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			chain, _ := test.chainFn() // sb
			stageProgressHandler := inmemory.NewStageProgressHandler(t)
			worker := v1.NewSparkTestWorker(t, chain, v1.WithIOHandler(inmemory.NewIOHandler(t)), v1.WithStageProgressHandler(stageProgressHandler))
			err := worker.Execute(context.NewJobMetadata(ctx.Background(), "jobKey", "correlationId", "transactionId", nil))
			if err != nil {
				t.Error(err)
			}
		})
	}
}
