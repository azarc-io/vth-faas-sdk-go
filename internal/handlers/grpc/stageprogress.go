package grpc

import (
	"context"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
)

type StageProgressHandler struct {
	client sdk_v1.ManagerServiceClient
}

func NewStageProgressHandler(client sdk_v1.ManagerServiceClient) api.StageProgressHandler {
	return &StageProgressHandler{client: client}
}

func (g *StageProgressHandler) Get(jobKey, name string) (*sdk_v1.StageStatus, error) {
	resp, err := g.client.GetStageStatus(context.Background(), sdk_v1.NewGetStageStatusReq(jobKey, name))
	return &resp.Status, err
}

func (g *StageProgressHandler) Set(stageStatus *sdk_v1.SetStageStatusRequest) error { //TODO receive ctx
	_, err := g.client.SetStageStatus(context.Background(), stageStatus)
	return err
}

func (g *StageProgressHandler) GetResult(jobKey, name string) (*sdk_v1.StageResult, error) { //TODO receive ctx
	result, err := g.client.GetStageResult(context.Background(), sdk_v1.NewStageResultReq(jobKey, name))
	return result.Result, err
}

func (g *StageProgressHandler) SetResult(result *sdk_v1.SetStageResultRequest) error { //TODO receive ctx
	_, err := g.client.SetStageResult(context.Background(), result)
	return err
}

func (g *StageProgressHandler) SetJobStatus(jobStatus *sdk_v1.SetJobStatusRequest) error { //TODO receive ctx
	_, err := g.client.SetJobStatus(context.Background(), jobStatus)
	return err
}

//func (g *StageProgressHandler) initialize() error {
//	if err := g.validate(); err != nil {
//		return err
//	}
//	if err := g.createClient(); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (g *StageProgressHandler) validate() error {
//	if g.config == nil {
//		return errors.New("a config is required in grpc stage progress handler")
//	}
//	return nil
//}
//
//func (g *StageProgressHandler) createClient() error {
//	var err error
//	g.client, err = CreateManagerServiceClient(g.config)
//	return err
//}

//type Option func(handler StageProgressHandler) StageProgressHandler
//
//func WithConfig(config *config.Config) Option {
//	return func(handler StageProgressHandler) StageProgressHandler {
//		handler.config = config
//		return handler
//	}
//}
