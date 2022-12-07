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

	err = worker.Execute(worker.LocalContext("test", "cid", "tid"))
	assert.Nil(t, err)

	var result string
	assert.Nil(t, io.Input("test", "message").Bind(&result))

	assert.Equal(t, "hello world", result)
	sph.AssertStageCompleted("test", "stage-1")
	sph.AssertStageCompleted("test", "stage-2")
	sph.AssertStageCompleted("test", "chain-1_complete")
}
