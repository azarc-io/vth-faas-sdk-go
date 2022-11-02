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

package grpc

import (
	"context"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
)

type StageProgressHandler struct {
	client sdk_v1.ManagerServiceClient
}

func NewStageProgressHandler(client sdk_v1.ManagerServiceClient) sdk_v1.StageProgressHandler {
	return &StageProgressHandler{client: client}
}

func (g *StageProgressHandler) Get(jobKey, name string) (*sdk_v1.StageStatus, error) {
	resp, err := g.client.GetStageStatus(context.Background(), sdk_v1.NewGetStageStatusReq(jobKey, name))
	return &resp.Status, err
}

func (g *StageProgressHandler) Set(stageStatus *sdk_v1.SetStageStatusRequest) error {
	_, err := g.client.SetStageStatus(context.Background(), stageStatus)
	return err
}

func (g *StageProgressHandler) GetResult(jobKey, name string) *sdk_v1.Result {
	result, err := g.client.GetStageResult(context.Background(), sdk_v1.NewStageResultReq(jobKey, name))
	if err != nil {
		return sdk_v1.NewResult(err, nil)
	}
	return sdk_v1.NewResult(nil, result.Result)
}

func (g *StageProgressHandler) SetResult(result *sdk_v1.SetStageResultRequest) error {
	_, err := g.client.SetStageResult(context.Background(), result)
	return err
}

func (g *StageProgressHandler) SetJobStatus(jobStatus *sdk_v1.SetJobStatusRequest) error {
	_, err := g.client.SetJobStatus(context.Background(), jobStatus)
	return err
}
