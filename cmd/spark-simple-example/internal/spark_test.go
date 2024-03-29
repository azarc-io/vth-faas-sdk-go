package spark_test

import (
	"context"
	spark "github.com/azarc-io/vth-faas-sdk-go/cmd/spark-simple-example/internal"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Should_Say_Hello_World(t *testing.T) {
	ctx := module_test_runner.NewTestJobContext(context.Background(),
		"say-hello-world", "cid", "tid", module_test_runner.Inputs{
			"myKey": {
				Value:    nil,
				MimeType: "",
			},
		})

	worker, err := module_test_runner.NewTestRunner(t, spark.NewSpark())
	assert.Nil(t, err)

	result, err := worker.Execute(
		ctx,
		sparkv1.WithSparkConfig(spark.Config{Foo: "my-bar-from-config"}))
	if !assert.Nil(t, err) {
		return
	}

	var message string
	assert.NoError(t, result.Bind("message", &message))
	assert.Equal(t, "hello world with bytes", message)

	worker.AssertStageCompleted("stage-1")
	worker.AssertStageCompleted("stage-2")
	worker.AssertStageCompleted("stage-3")
	worker.AssertStageCompleted("stage-5")
	worker.AssertStageCompleted("chain-1_complete")
	worker.AssertStageOrder("stage-1", "stage-2", "stage-3", "stage-4", "stage-5")
	worker.AssertStageResult("stage-4", "my-bar-from-config")
	worker.AssertStageResult("stage-5", "JobKey:say-hello-world; TransactionId:tid; CorrelationId:cid")
}

func Test_Should_Cancel(t *testing.T) {
	bCtx, cancel := context.WithCancel(context.Background())
	ctx := module_test_runner.NewTestJobContext(bCtx, "should_cancel", "cid", "tid", module_test_runner.Inputs{
		"myKey": {
			Value:    nil,
			MimeType: "",
		},
	})

	worker, err := module_test_runner.NewTestRunner(t, spark.NewSpark())
	assert.Nil(t, err)

	go func() {
		// force cancel after 500ms
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()

	result, err := worker.Execute(ctx, sparkv1.WithSparkConfig(spark.Config{Foo: "my-bar-from-config"}))
	if !assert.ErrorContains(t, err, "canceled") {
		return
	}

	assert.Nil(t, result)

	worker.AssertStageCompleted("stage-1")
	worker.AssertStageCompleted("stage-2")
	worker.AssertStageFailed("stage-3")
	worker.AssertStageOrder("stage-1", "stage-2")
}
