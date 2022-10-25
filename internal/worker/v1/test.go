package v1

import (
	"github.com/azarc-io/vth-faas-sdk-go/internal/spark"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/config"
	"testing"
)

type TestWorker struct {
	t      *testing.T
	worker sdk_v1.Worker
}

func (t TestWorker) Execute(ctx sdk_v1.Context) sdk_v1.StageError {
	return t.worker.Execute(ctx)
}

func NewSparkTestWorker(t *testing.T, chain *spark.Chain, options ...Option) sdk_v1.Worker {
	cfg, err := config.NewMock(map[string]string{"APP_ENVIRONMENT": "test", "AGENT_SERVER_PORT": "0", "MANAGER_SERVER_PORT": "0"})
	if err != nil {
		t.Error(err)
	}
	sw, err := NewSparkWorker(cfg, chain, options...)
	if err != nil {
		t.Fatal(err)
		return nil
	}
	return &TestWorker{
		t:      t,
		worker: sw,
	}
}
