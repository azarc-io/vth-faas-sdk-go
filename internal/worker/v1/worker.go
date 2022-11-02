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

package v1

import (
	"github.com/azarc-io/vth-faas-sdk-go/internal/clients"
	"github.com/azarc-io/vth-faas-sdk-go/internal/context"
	grpc_handler "github.com/azarc-io/vth-faas-sdk-go/internal/handlers/grpc"
	"github.com/azarc-io/vth-faas-sdk-go/internal/logger"
	"github.com/azarc-io/vth-faas-sdk-go/internal/spark"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/config"
)

type SparkWorker struct {
	config               *config.Config
	chain                *spark.Chain
	variableHandler      sdk_v1.IOHandler
	stageProgressHandler sdk_v1.StageProgressHandler
	log                  sdk_v1.Logger
}

func NewSparkWorker(cfg *config.Config, chain *spark.Chain, options ...Option) (sdk_v1.Worker, error) {
	jw := &SparkWorker{config: cfg, chain: chain}
	for _, opt := range options {
		jw = opt(jw)
	}
	if err := jw.validate(); err != nil {
		return nil, err
	}
	return jw, nil
}

func (w *SparkWorker) Execute(metadata sdk_v1.Context) sdk_v1.StageError {
	jobContext := context.NewJobContext(metadata, w.stageProgressHandler, w.variableHandler, w.log)
	return w.chain.Execute(jobContext)
}

func (w *SparkWorker) validate() error {
	var grpcClient sdk_v1.ManagerServiceClient
	if w.variableHandler == nil || w.stageProgressHandler == nil {
		var err error
		grpcClient, err = clients.CreateManagerServiceClient(w.config)
		if err != nil {
			return err
		}
	}
	if w.variableHandler == nil {
		w.variableHandler = grpc_handler.NewIOHandler(grpcClient)
	}
	if w.stageProgressHandler == nil {
		w.stageProgressHandler = grpc_handler.NewStageProgressHandler(grpcClient)
	}
	if w.log == nil {
		w.log = logger.NewLogger()
	}
	return nil
}

type Option = func(je *SparkWorker) *SparkWorker

func WithIOHandler(vh sdk_v1.IOHandler) Option {
	return func(jw *SparkWorker) *SparkWorker {
		jw.variableHandler = vh
		return jw
	}
}

func WithStageProgressHandler(sph sdk_v1.StageProgressHandler) Option {
	return func(jw *SparkWorker) *SparkWorker {
		jw.stageProgressHandler = sph
		return jw
	}
}

func WithLog(log sdk_v1.Logger) Option {
	return func(jw *SparkWorker) *SparkWorker {
		jw.log = log
		return jw
	}
}
