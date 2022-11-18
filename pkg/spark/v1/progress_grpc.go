package spark_v1

import (
	"context"
)

type stageProgressHandler struct {
	client ManagerServiceClient
}

func newGrpcStageProgressHandler(client ManagerServiceClient) StageProgressHandler {
	return &stageProgressHandler{client: client}
}

func (g *stageProgressHandler) Get(jobKey, name string) (*StageStatus, error) {
	resp, err := g.client.GetStageStatus(context.Background(), newGetStageStatusReq(jobKey, name))
	return &resp.Status, err
}

func (g *stageProgressHandler) Set(stageStatus *SetStageStatusRequest) error {
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

func (g *stageProgressHandler) SetResult(result *SetStageResultRequest) error {
	_, err := g.client.SetStageResult(context.Background(), result)
	return err
}

func (g *stageProgressHandler) SetJobStatus(jobStatus *SetJobStatusRequest) error {
	_, err := g.client.SetJobStatus(context.Background(), jobStatus)
	return err
}
