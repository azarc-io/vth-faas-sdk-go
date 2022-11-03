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
