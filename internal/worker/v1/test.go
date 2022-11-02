// Copyright 2020-2022 Azarc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1 // TODO do not add that to the binary

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
