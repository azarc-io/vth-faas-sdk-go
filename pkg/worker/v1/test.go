package v1 // TODO do not add that to the binary

import (
	"testing"

	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1"

	"github.com/azarc-io/vth-faas-sdk-go/pkg/config"
)

type TestWorker struct {
	t      *testing.T
	worker sdk_v1.Worker
}

func (t TestWorker) Execute(ctx sdk_v1.Context) sdk_v1.StageError {
	return t.worker.Execute(ctx)
}

func NewSparkTestWorker(t *testing.T, chain *sdk_v1.BuilderChain, options ...sdk_v1.ErrorOption) sdk_v1.Worker {
	cfg, err := config.NewMock(map[string]string{"APP_ENVIRONMENT": "test", "AGENT_SERVER_PORT": "0", "MANAGER_SERVER_PORT": "0"})
	if err != nil {
		t.Error(err)
	}
	sw, err := sdk_v1.NewSparkWorker(cfg, chain, options...)
	if err != nil {
		t.Fatal(err)
		return nil
	}
	return &TestWorker{
		t:      t,
		worker: sw,
	}
}
