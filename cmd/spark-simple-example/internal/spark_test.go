package spark_test

import (
	"context"
	spark "github.com/azarc-io/vth-faas-sdk-go/cmd/spark-simple-example/internal"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Should_Say_Hello_World(t *testing.T) {
	io := sparkv1.NewInMemoryIOHandler(t)
	sph := sparkv1.NewInMemoryStageProgressHandler(t)

	worker, err := sparkv1.NewSparkWorker(
		context.Background(),
		spark.NewSpark(),
		sparkv1.WithIOHandler(io),
		sparkv1.WithStageProgressHandler(sph),
	)
	assert.Nil(t, err)

	ctx := worker.LocalContext("test", "cid", "tid")
	err = worker.Execute(ctx)
	assert.Nil(t, err)

	var result string
	assert.Nil(t, io.Input(ctx, "message").Bind(&result))

	assert.Equal(t, "hello world", result)
	sph.AssertStageCompleted(ctx, "stage-1")
	sph.AssertStageCompleted(ctx, "stage-2")
	sph.AssertStageCompleted(ctx, "chain-1_complete")
}
