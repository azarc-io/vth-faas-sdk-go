package spark_test

import (
	"context"
	spark "github.com/azarc-io/vth-faas-sdk-go/cmd/simple-example/internal"
	spark_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Should_Say_Hello_World(t *testing.T) {
	io := spark_v1.NewInMemoryIOHandler(t)
	sph := spark_v1.NewInMemoryStageProgressHandler(t)

	worker, err := spark_v1.NewSparkWorker(
		context.Background(),
		spark.NewSpark(),
		spark_v1.WithIOHandler(io),
		spark_v1.WithStageProgressHandler(sph),
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
