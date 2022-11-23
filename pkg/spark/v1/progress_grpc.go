package spark_v1

import (
	"context"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/internal/gen/azarc/sdk/spark/v1"
)

type stageProgressHandler struct {
	client sparkv1.ManagerServiceClient
}

func newGrpcStageProgressHandler(client sparkv1.ManagerServiceClient) StageProgressHandler {
	return &stageProgressHandler{client: client}
}

func (g *stageProgressHandler) Get(jobKey, name string) (*sparkv1.StageStatus, error) {
	resp, err := g.client.GetStageStatus(context.Background(), newGetStageStatusReq(jobKey, name))
	return &resp.Status, err
}

func (g *stageProgressHandler) Set(stageStatus *sparkv1.SetStageStatusRequest) error {
	_, err := g.client.SetStageStatus(context.Background(), stageStatus)
	return err
}

func (g *stageProgressHandler) GetResult(jobKey, name string) Bindable {
	result, err := g.client.GetStageResult(context.Background(), newStageResultReq(jobKey, name))
	if err != nil {
		return newResult(err, nil)
	}
	return newResult(nil, result.Result)
}

func (g *stageProgressHandler) SetResult(result *sparkv1.SetStageResultRequest) error {
	_, err := g.client.SetStageResult(context.Background(), result)
	return err
}

func (g *stageProgressHandler) SetJobStatus(jobStatus *sparkv1.SetJobStatusRequest) error {
	_, err := g.client.SetJobStatus(context.Background(), jobStatus)
	return err
}
